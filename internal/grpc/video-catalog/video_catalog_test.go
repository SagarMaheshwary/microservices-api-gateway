package videocatalog_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	videocatalog "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video-catalog"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

var dummyVideo = &videocatalogpb.Video{
	Id:           1,
	Title:        "title",
	Description:  "description",
	ThumbnailUrl: "example.com",
	PublishedAt:  time.Now().String(),
	Duration:     100,
	Resolution:   "1080x1920",
	User: &videocatalogpb.User{
		Id:    1,
		Name:  "name",
		Image: nil,
	},
}

func TestVideoCatalogClient_FindAll(t *testing.T) {
	req := &videocatalogpb.FindAllRequest{}
	res := &videocatalogpb.FindAllResponse{
		Message: constant.MessageOK,
		Data: &videocatalogpb.FindAllResponseData{
			Videos: []*videocatalogpb.Video{dummyVideo},
		},
	}

	cfg := &config.GRPCVideoCatalogClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *videocatalogpb.FindAllResponse
		mockErr    error
		expectErr  bool
	}{
		{
			name:       "success",
			mockReturn: res,
			mockErr:    nil,
			expectErr:  false,
		},
		{
			name:       "gRPC error",
			mockReturn: nil,
			mockErr:    errors.New("grpc error"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockVideoCatalogServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("FindAll", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := videocatalog.NewVideoCatalogClient(mockClient, mockHealth, cfg)

			got, err := c.FindAll(context.Background(), req)

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestVideoCatalogClient_FindById(t *testing.T) {
	req := &videocatalogpb.FindByIdRequest{}
	res := &videocatalogpb.FindByIdResponse{
		Message: constant.MessageOK,
		Data: &videocatalogpb.FindByIdResponseData{
			Video: dummyVideo,
		},
	}

	cfg := &config.GRPCVideoCatalogClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *videocatalogpb.FindByIdResponse
		mockErr    error
		expectErr  bool
		expectGRPC codes.Code
	}{
		{
			name:       "success",
			mockReturn: res,
			mockErr:    nil,
			expectErr:  false,
		},
		{
			name:       "gRPC error",
			mockReturn: nil,
			mockErr:    errors.New("grpc error"),
			expectErr:  true,
		},
		{
			name:       "not found",
			mockReturn: nil,
			mockErr:    status.Error(codes.NotFound, "video not found"),
			expectErr:  true,
			expectGRPC: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockVideoCatalogServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("FindById", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := videocatalog.NewVideoCatalogClient(mockClient, mockHealth, cfg)

			got, err := c.FindById(context.Background(), req)

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, got)

				// if we expect a specific gRPC code
				if tt.expectGRPC != 0 {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.expectGRPC, st.Code())
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, got)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestUploadClient_Health(t *testing.T) {
	cfg := &config.GRPCVideoCatalogClient{Timeout: 2 * time.Second}

	tests := []struct {
		name      string
		mockResp  *healthpb.HealthCheckResponse
		mockErr   error
		expectErr bool
		expectMsg string
	}{
		{
			name:      "success",
			mockResp:  &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING},
			mockErr:   nil,
			expectErr: false,
		},
		{
			name:      "gRPC error",
			mockResp:  nil,
			mockErr:   errors.New("grpc health check failed"),
			expectErr: true,
			expectMsg: "grpc health check failed",
		},
		{
			name:      "not serving",
			mockResp:  &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING},
			mockErr:   nil,
			expectErr: true,
			expectMsg: "video catalog grpc health check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpload := new(MockVideoCatalogServiceClient)
			mockHealth := new(MockHealthClient)

			client := videocatalog.NewVideoCatalogClient(mockUpload, mockHealth, cfg)

			mockHealth.On("Check", mock.Anything, &healthpb.HealthCheckRequest{}).
				Return(tt.mockResp, tt.mockErr).Once()

			err := client.Health(context.Background())

			if tt.expectErr {
				require.Error(t, err)
				if tt.expectMsg != "" {
					assert.EqualError(t, err, tt.expectMsg)
				}
			} else {
				require.NoError(t, err)
			}

			mockHealth.AssertExpectations(t)
		})
	}
}
