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

type DialFunc func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

type ClientFactory func(c authpb.AuthenticationServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) AuthenticationService

type InitClientOptions struct {
	Config          *config.GRPCClient
	Dial            DialFunc
	Factory         ClientFactory
	DialOptions     []grpc.DialOption
	SkipHealthCheck bool
}

func defaultDialer(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

func defaultFactory(c authpb.AuthenticationServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) AuthenticationService {
	return NewAuthenticationClient(c, h, cfg)
}

func NewClient(ctx context.Context, opt *InitClientOptions) (AuthenticationService, *grpc.ClientConn, error) {
	if opt == nil {
		opt = &InitClientOptions{}
	}
	if opt.Dial == nil {
		opt.Dial = defaultDialer
	}
	if opt.Factory == nil {
		opt.Factory = defaultFactory
	}
	if len(opt.DialOptions) == 0 {
		opt.DialOptions = []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithStatsHandler(otelgrpc.NewClientHandler(
				otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
				otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
			)),
		}
	}

	conn, err := opt.Dial(opt.Config.UploadServiceURL, opt.DialOptions...)
	if err != nil {
		logger.Error("Authentication gRPC client failed to connect on %q: %v", opt.Config.UploadServiceURL, err)
		return nil, nil, err
	}

	logger.Info("Authentication gRPC client connected on %q", opt.Config.UploadServiceURL)

	uploadClient := opt.Factory(
		authpb.NewAuthenticationServiceClient(conn),
		healthpb.NewHealthClient(conn),
		opt.Config,
	)

	if !opt.SkipHealthCheck {
		if err := uploadClient.Health(ctx); err != nil {
			return nil, nil, err
		}
	}

	logger.Info("Authentication gRPC client ready on %q", opt.Config.UploadServiceURL)
	return uploadClient, conn, nil
}
