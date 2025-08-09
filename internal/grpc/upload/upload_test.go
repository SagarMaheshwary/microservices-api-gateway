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
	mockClient := new(MockUploadServiceClient)
	mockHealth := new(MockHealthClient)

	cfg := &config.GRPCClient{Timeout: 2 * time.Second}
	c := upload.NewUploadClient(mockClient, mockHealth, cfg)

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

	// Success case
	mockClient.On("CreatePresignedUrl", mock.Anything, req).Return(res, nil)
	got, err := c.CreatePresignedUrl(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, res, got)

	// Failure case
	mockClient.On("CreatePresignedUrl", mock.Anything, req).Return(nil, errors.New("grpc failure"))
	got, err = c.CreatePresignedUrl(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, res, got)

	mockClient.AssertExpectations(t)
}

func TestUploadClient_UploadedWebhook(t *testing.T) {
	mockUpload := new(MockUploadServiceClient)
	mockHealth := new(MockHealthClient)

	cfg := &config.GRPCClient{Timeout: 2 * time.Second}
	client := upload.NewUploadClient(mockUpload, mockHealth, cfg)

	req := &uploadpb.UploadedWebhookRequest{
		VideoId:     "1",
		ThumbnailId: "1",
		Title:       "title",
		Description: "description",
	}
	resp := &uploadpb.UploadedWebhookResponse{Message: constant.MessageOK, Data: &uploadpb.UploadedWebhookResponseData{}}
	userID := "user-1"

	// Success case
	mockUpload.On("UploadedWebhook", mock.Anything, req).Return(resp, nil).Once()
	got, err := client.UploadedWebhook(context.Background(), req, userID)
	require.NoError(t, err)
	assert.Equal(t, resp, got)

	// Failure case
	mockUpload.On("UploadedWebhook", mock.Anything, req).Return(nil, errors.New("grpc failure")).Once()
	got, err = client.UploadedWebhook(context.Background(), req, userID)
	require.Error(t, err)
	assert.Nil(t, got)

	mockUpload.AssertExpectations(t)
}

func TestUploadClient_Health(t *testing.T) {
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

			cfg := &config.GRPCClient{Timeout: 2 * time.Second}
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
