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

func (u *uploadClient) CreatePresignedUrl(data *uploadpb.CreatePresignedUrlRequest) (*uploadpb.CreatePresignedUrlResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := u.client.CreatePresignedUrl(ctx, data)

	if err != nil {
		logger.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)

		return nil, err
	}

	return response, nil
}

func (u *uploadClient) UploadedWebhook(data *uploadpb.UploadedWebhookRequest, userId string) (*uploadpb.UploadedWebhookResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	md := metadata.Pairs(constant.HeaderUserId, userId)
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	response, err := u.client.UploadedWebhook(ctx, data)

	if err != nil {
		logger.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)

		return nil, err
	}

	return response, nil
}
