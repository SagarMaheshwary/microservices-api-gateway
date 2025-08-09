package upload_test

import (
	"context"

	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type MockUploadServiceClient struct {
	mock.Mock
}

func (m *MockUploadServiceClient) CreatePresignedUrl(ctx context.Context, in *uploadpb.CreatePresignedUrlRequest, opts ...grpc.CallOption) (*uploadpb.CreatePresignedUrlResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*uploadpb.CreatePresignedUrlResponse), nil
}

func (m *MockUploadServiceClient) UploadedWebhook(ctx context.Context, in *uploadpb.UploadedWebhookRequest, opts ...grpc.CallOption) (*uploadpb.UploadedWebhookResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*uploadpb.UploadedWebhookResponse), nil
}

func (m *MockUploadServiceClient) Health(ctx context.Context, opts ...grpc.CallOption) error {
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
