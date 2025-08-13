package handler_test

import (
	"context"

	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"github.com/stretchr/testify/mock"
)

type MockUploadServiceClient struct {
	mock.Mock
}

func (m *MockUploadServiceClient) CreatePresignedUrl(ctx context.Context, in *uploadpb.CreatePresignedUrlRequest) (*uploadpb.CreatePresignedUrlResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*uploadpb.CreatePresignedUrlResponse), nil
}

func (m *MockUploadServiceClient) UploadedWebhook(ctx context.Context, in *uploadpb.UploadedWebhookRequest, token string) (*uploadpb.UploadedWebhookResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*uploadpb.UploadedWebhookResponse), nil
}

func (m *MockUploadServiceClient) Health(ctx context.Context) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}
