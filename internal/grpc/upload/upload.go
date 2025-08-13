package upload

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

type UploadService interface {
	CreatePresignedUrl(ctx context.Context, in *uploadpb.CreatePresignedUrlRequest) (*uploadpb.CreatePresignedUrlResponse, error)
	UploadedWebhook(ctx context.Context, data *uploadpb.UploadedWebhookRequest, userId string) (*uploadpb.UploadedWebhookResponse, error)
	Health(ctx context.Context) error
}

type UploadClient struct {
	config *config.GRPCUploadClient
	client uploadpb.UploadServiceClient
	health healthpb.HealthClient
}

func NewUploadClient(c uploadpb.UploadServiceClient, h healthpb.HealthClient, cfg *config.GRPCUploadClient) *UploadClient {
	return &UploadClient{client: c, health: h, config: cfg}
}

func (u *UploadClient) CreatePresignedUrl(ctx context.Context, in *uploadpb.CreatePresignedUrlRequest) (*uploadpb.CreatePresignedUrlResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	response, err := u.client.CreatePresignedUrl(ctx, in)
	if err != nil {
		logger.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)
		return nil, err
	}

	return response, nil
}

func (u *UploadClient) UploadedWebhook(ctx context.Context, in *uploadpb.UploadedWebhookRequest, userId string) (*uploadpb.UploadedWebhookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	md := metadata.Pairs(constant.GRPCHeaderUserId, userId)
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := u.client.UploadedWebhook(ctx, in)
	if err != nil {
		logger.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)
		return nil, err
	}

	return response, nil
}

func (u *UploadClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, u.config.Timeout)
	defer cancel()

	res, err := u.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		logger.Error("Upload gRPC health check failed! %v", err)
		return err
	}

	if res.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("Upload gRPC health check failed")
		return errors.New("upload grpc health check failed")
	}

	return nil
}
