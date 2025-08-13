package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_CheckAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(auth *MockAuthenticationServiceClient, upload *MockUploadServiceClient, video *MockVideoCatalogServiceClient)
		expectedStatus int
		expectedJSON   gin.H
	}{
		{
			name: "success",
			mockSetup: func(auth *MockAuthenticationServiceClient, upload *MockUploadServiceClient, video *MockVideoCatalogServiceClient) {
				auth.On("Health", mock.Anything).
					Return(nil).
					Once()

				upload.On("Health", mock.Anything).
					Return(nil).
					Once()

				video.On("Health", mock.Anything).
					Return(nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageServicesHealthy,
				"data": gin.H{
					"status": constant.HealthStatusHealthy,
				},
			},
		},
		{
			name: "unhealthy/grpc error",
			mockSetup: func(auth *MockAuthenticationServiceClient, upload *MockUploadServiceClient, video *MockVideoCatalogServiceClient) {
				// Only mock the first one the handler will check.
				auth.On("Health", mock.Anything).
					Return(errors.New("grpc failed")).
					Once()
				// No mocks for upload/video because they won’t be called.
			},
			expectedStatus: http.StatusServiceUnavailable,
			expectedJSON: gin.H{
				"message": constant.MessageServicesUnhealthy,
				"data": gin.H{
					"status": constant.HealthStatusDegraded,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthenticationClient := new(MockAuthenticationServiceClient)
			mockUploadClient := new(MockUploadServiceClient)
			mockVideoCatalogSvc := new(MockVideoCatalogServiceClient)

			tt.mockSetup(mockAuthenticationClient, mockUploadClient, mockVideoCatalogSvc)

			h := handler.NewHealthHandler(types.GRPCClients{
				AuthClient:         mockAuthenticationClient,
				UploadClient:       mockUploadClient,
				VideoCatalogClient: mockVideoCatalogSvc,
			})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)
			c.Request.Header.Set("Content-Type", "application/json")

			h.CheckAll(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockAuthenticationClient.AssertExpectations(t)
			mockUploadClient.AssertExpectations(t)
			mockVideoCatalogSvc.AssertExpectations(t)
		})
	}
}
