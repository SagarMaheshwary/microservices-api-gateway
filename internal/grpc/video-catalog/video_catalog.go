package videocatalog

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type VideoCatalogService interface {
	FindAll(ctx context.Context, in *videocatalogpb.FindAllRequest) (*videocatalogpb.FindAllResponse, error)
	FindById(ctx context.Context, in *videocatalogpb.FindByIdRequest) (*videocatalogpb.FindByIdResponse, error)
	Health(ctx context.Context) error
}

type VideoCatalogClient struct {
	config *config.GRPCClient
	client videocatalogpb.VideoCatalogServiceClient
	health healthpb.HealthClient
}

func NewVideoCatalogClient(c videocatalogpb.VideoCatalogServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) *VideoCatalogClient {
	return &VideoCatalogClient{
		client: c,
		health: h,
		config: cfg,
	}
}

func (v *VideoCatalogClient) FindAll(ctx context.Context, in *videocatalogpb.FindAllRequest) (*videocatalogpb.FindAllResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, v.config.Timeout)
	defer cancel()

	response, err := v.client.FindAll(ctx, in)
	if err != nil {
		logger.Error("gRPC videoCatalogClient.FindAll failed %v", err)
		return nil, err
	}

	return response, nil
}

func (v *VideoCatalogClient) FindById(ctx context.Context, in *videocatalogpb.FindByIdRequest) (*videocatalogpb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, v.config.Timeout)
	defer cancel()

	response, err := v.client.FindById(ctx, in)
	if err != nil {
		logger.Error("gRPC videoCatalogClient.FindById failed %v", err)
		return nil, err
	}

	return response, nil
}

func (v *VideoCatalogClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, v.config.Timeout)
	defer cancel()

	res, err := v.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		logger.Error("Video Catalog gRPC health check failed! %v", err)
		return err
	}

	if res.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("Video Catalog gRPC health check failed")
		return errors.New("video catalog grpc health check failed")
	}

	return nil
}
