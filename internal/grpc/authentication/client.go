package authentication

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func InitClient(ctx context.Context) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	opts = append(
		opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
	)

	address := config.Conf.GRPCClient.AuthenticationServiceURL

	conn, err := grpc.NewClient(address, opts...)

	if err != nil {
		logger.Error("gRPC client failed to connect on %q: %v", address, err)
		return nil, err
	}

	logger.Info("gRPC client connected on %q", address)

	Auth = &authenticationClient{
		client: authpb.NewAuthenticationServiceClient(conn),
		health: healthpb.NewHealthClient(conn),
	}

	if err := HealthCheck(ctx); err != nil {
		return nil, err
	}

	logger.Info("Authentication gRPC client connected on %q", address)

	return conn, nil
}

func HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.TimeoutSeconds)
	defer cancel()

	response, err := Auth.health.Check(ctx, &healthpb.HealthCheckRequest{})

	if err != nil {
		logger.Error("Authentication gRPC health check failed! %v", err)
		return err
	}

	if response.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("Authentication gRPC health check failed!")
		return errors.New("Authentication gRPC health check failed")
	}

	return nil
}
