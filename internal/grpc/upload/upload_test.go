package upload_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestUploadClient_CreatePresignedUrl(t *testing.T) {
	req := &uploadpb.CreatePresignedUrlRequest{}
	res := &uploadpb.CreatePresignedUrlResponse{
		Message: constant.MessageOK,
		Data: &uploadpb.CreatePresignedUrlResponseData{
			ThumbnailUrl: "https://example.com/presigned",
			VideoId:      "1",
			ThumbnailId:  "1",
			VideoUrl:     "https://example.com/presigned",
		},
	}
	cfg := &config.GRPCUploadClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *uploadpb.CreatePresignedUrlResponse
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
			name:       "generic failure",
			mockReturn: nil,
			mockErr:    errors.New("grpc failure"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUploadServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("CreatePresignedUrl", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := upload.NewUploadClient(mockClient, mockHealth, cfg)

			got, err := c.CreatePresignedUrl(context.Background(), req)

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

func TestUploadClient_UploadedWebhook(t *testing.T) {
	req := &uploadpb.UploadedWebhookRequest{
		VideoId:     "1",
		ThumbnailId: "1",
		Title:       "title",
		Description: "description",
	}
	res := &uploadpb.UploadedWebhookResponse{Message: constant.MessageOK, Data: &uploadpb.UploadedWebhookResponseData{}}

	cfg := &config.GRPCUploadClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *uploadpb.UploadedWebhookResponse
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
			name:       "generic failure",
			mockReturn: nil,
			mockErr:    errors.New("grpc failure"),
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUploadServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("UploadedWebhook", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := upload.NewUploadClient(mockClient, mockHealth, cfg)

			got, err := c.UploadedWebhook(context.Background(), req, "user-id")

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

func TestUploadClient_Health(t *testing.T) {
	cfg := &config.GRPCUploadClient{Timeout: 2 * time.Second}

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
			expectMsg: "upload grpc health check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpload := new(MockUploadServiceClient)
			mockHealth := new(MockHealthClient)

			client := upload.NewUploadClient(mockUpload, mockHealth, cfg)

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
