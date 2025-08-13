package handler_test

import (
	"context"

	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"github.com/stretchr/testify/mock"
)

type MockVideoCatalogServiceClient struct {
	mock.Mock
}

func (m *MockVideoCatalogServiceClient) FindAll(ctx context.Context, in *videocatalogpb.FindAllRequest) (*videocatalogpb.FindAllResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*videocatalogpb.FindAllResponse), nil
}

func (m *MockVideoCatalogServiceClient) FindById(ctx context.Context, in *videocatalogpb.FindByIdRequest) (*videocatalogpb.FindByIdResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*videocatalogpb.FindByIdResponse), nil
}

func (m *MockVideoCatalogServiceClient) Health(ctx context.Context) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}
