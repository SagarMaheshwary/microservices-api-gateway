package upload

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

var Upload *uploadClient

type uploadClient struct {
	client uploadpb.UploadServiceClient
	health healthpb.HealthClient
}

func (u *uploadClient) CreatePresignedUrl(ctx context.Context, data *uploadpb.CreatePresignedUrlRequest) (*uploadpb.CreatePresignedUrlResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := u.client.CreatePresignedUrl(ctx, data)
	if err != nil {
		logger.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)
		return nil, err
	}

	return response, nil
}

func (u *uploadClient) UploadedWebhook(ctx context.Context, data *uploadpb.UploadedWebhookRequest, userId string) (*uploadpb.UploadedWebhookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	md := metadata.Pairs(constant.GRPCHeaderUserId, userId)
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := u.client.UploadedWebhook(ctx, data)
	if err != nil {
		logger.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)
		return nil, err
	}

	return response, nil
}
