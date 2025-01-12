package video_catalog

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var VideoCatalog *videoCatalogClient

type videoCatalogClient struct {
	client videocatalogpb.VideoCatalogServiceClient
	health healthpb.HealthClient
}

func (v *videoCatalogClient) FindAll(data *videocatalogpb.FindAllRequest) (*videocatalogpb.FindAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := v.client.FindAll(ctx, data)

	if err != nil {
		logger.Error("gRPC videoCatalogClient.FindAll failed %v", err)

		return nil, err
	}

	return response, nil
}

func (v *videoCatalogClient) FindById(data *videocatalogpb.FindByIdRequest) (*videocatalogpb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := v.client.FindById(ctx, data)

	if err != nil {
		logger.Error("gRPC videoCatalogClient.FindById failed %v", err)

		return nil, err
	}

	return response, nil
}
