package authentication

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func Connect() {
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

	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		logger.Error("gRPC client failed to connect on %q: %v", address, err)
	}

	logger.Info("gRPC client connected on %q", address)

	Auth = &authenticationClient{
		client: authpb.NewAuthenticationServiceClient(conn),
		health: healthpb.NewHealthClient(conn),
	}
}

func HealthCheck(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, config.Conf.GRPCClient.Timeout)
	defer cancel()

	response, err := Auth.health.Check(ctx, &healthpb.HealthCheckRequest{})

	if err != nil {
		logger.Error("Authentication gRPC health check failed! %v", err)

		return false
	}

	if response.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("Authentication gRPC health check failed!")

		return false
	}

	return true
}
