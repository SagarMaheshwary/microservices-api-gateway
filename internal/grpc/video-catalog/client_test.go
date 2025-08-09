package videocatalog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	videocatalog "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video-catalog"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestInitClient_Success(t *testing.T) {
	mockClient := new(MockVideoCatalogServiceClient)
	mockHealth := new(MockHealthClient)

	opt := &videocatalog.InitClientOptions{
		Config: &config.GRPCClient{UploadServiceURL: "fake-url"},
		Dial: func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
			return &grpc.ClientConn{}, nil
		},
		Factory: func(c videocatalogpb.VideoCatalogServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) videocatalog.VideoCatalogService {
			return videocatalog.NewVideoCatalogClient(mockClient, mockHealth, cfg)
		},
	}

	mockHealth.On("Check", mock.Anything, mock.Anything, mock.Anything).Return(
		&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil,
	)

	svc, conn, err := videocatalog.NewClient(context.Background(), opt)

	assert.NoError(t, err)
	assert.NotNil(t, svc)
	assert.NotNil(t, conn)

	mockHealth.AssertExpectations(t)
}

func TestInitClient_HealthFails(t *testing.T) {
	mockClient := new(MockVideoCatalogServiceClient)
	mockHealth := new(MockHealthClient)

	opt := &videocatalog.InitClientOptions{
		Config: &config.GRPCClient{UploadServiceURL: "fake-url"},
		Dial: func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
			return &grpc.ClientConn{}, nil
		},
		Factory: func(c videocatalogpb.VideoCatalogServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) videocatalog.VideoCatalogService {
			return videocatalog.NewVideoCatalogClient(mockClient, mockHealth, cfg)
		},
	}

	mockHealth.On("Check", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, errors.New("health failed"),
	)

	svc, conn, err := videocatalog.NewClient(context.Background(), opt)

	assert.Error(t, err)
	assert.Nil(t, svc)
	assert.Nil(t, conn)

	mockHealth.AssertExpectations(t)
}

func TestInitClient_DialFail(t *testing.T) {
	dialErr := errors.New("dial failed")

	opt := &videocatalog.InitClientOptions{
		Config: &config.GRPCClient{UploadServiceURL: "fake-url"},
		Dial: func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
			return nil, dialErr
		},
		Factory: func(c videocatalogpb.VideoCatalogServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) videocatalog.VideoCatalogService {
			t.Fatal("Factory should not be called when dial fails")
			return nil
		},
	}

	svc, conn, err := videocatalog.NewClient(context.Background(), opt)

	assert.EqualError(t, err, dialErr.Error())
	assert.Nil(t, svc)
	assert.Nil(t, conn)
}
