package upload

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type DialFunc func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

type ClientFactory func(c uploadpb.UploadServiceClient, h healthpb.HealthClient, cfg *config.GRPCUploadClient) UploadService

type InitClientOptions struct {
	Config          *config.GRPCUploadClient
	Dial            DialFunc
	Factory         ClientFactory
	DialOptions     []grpc.DialOption
	SkipHealthCheck bool
}

func defaultDialer(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

func defaultFactory(c uploadpb.UploadServiceClient, h healthpb.HealthClient, cfg *config.GRPCUploadClient) UploadService {
	return NewUploadClient(c, h, cfg)
}

func NewClient(ctx context.Context, opt *InitClientOptions) (UploadService, *grpc.ClientConn, error) {
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

	conn, err := opt.Dial(opt.Config.URL, opt.DialOptions...)
	if err != nil {
		logger.Error("Upload gRPC client failed to connect on %q: %v", opt.Config.URL, err)
		return nil, nil, err
	}

	logger.Info("Upload gRPC client connected on %q", opt.Config.URL)

	uploadClient := opt.Factory(
		uploadpb.NewUploadServiceClient(conn),
		healthpb.NewHealthClient(conn),
		opt.Config,
	)

	if !opt.SkipHealthCheck {
		if err := uploadClient.Health(ctx); err != nil {
			return nil, nil, err
		}
	}

	logger.Info("Upload gRPC client ready on %q", opt.Config.URL)
	return uploadClient, conn, nil
}
