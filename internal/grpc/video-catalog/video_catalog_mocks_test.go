package videocatalog_test

import (
	"context"

	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type MockVideoCatalogServiceClient struct {
	mock.Mock
}

func (m *MockVideoCatalogServiceClient) FindAll(ctx context.Context, in *videocatalogpb.FindAllRequest, opts ...grpc.CallOption) (*videocatalogpb.FindAllResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*videocatalogpb.FindAllResponse), nil
}

func (m *MockVideoCatalogServiceClient) FindById(ctx context.Context, in *videocatalogpb.FindByIdRequest, opts ...grpc.CallOption) (*videocatalogpb.FindByIdResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*videocatalogpb.FindByIdResponse), nil
}

func (m *MockVideoCatalogServiceClient) Health(ctx context.Context, opts ...grpc.CallOption) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

type MockHealthClient struct {
	mock.Mock
	healthpb.HealthClient
}

func (m *MockHealthClient) Check(ctx context.Context, in *healthpb.HealthCheckRequest, opts ...grpc.CallOption) (*healthpb.HealthCheckResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*healthpb.HealthCheckResponse), args.Error(1)
}
