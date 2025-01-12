package video_catalog

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address := config.Conf.GRPCClient.VideoCatalogServiceurl

	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		logger.Error("gRPC client failed to connect on %q: %v", address, err)
	}

	logger.Info("gRPC client connected on %q", address)

	VideoCatalog = &videoCatalogClient{
		client: videocatalogpb.NewVideoCatalogServiceClient(conn),
		health: healthpb.NewHealthClient(conn),
	}
}

func HealthCheck() bool {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)
	defer cancel()

	response, err := VideoCatalog.health.Check(ctx, &healthpb.HealthCheckRequest{})

	if err != nil {
		logger.Error("VideoCatalog gRPC health check failed! %v", err)

		return false
	}

	if response.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("VideoCatalog gRPC health check failed!")

		return false
	}

	return true
}
