package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUploadHandler_CreatePresignedUrl(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(m *MockUploadServiceClient)
		expectedStatus int
		expectedJSON   gin.H
		mockVerifyFunc middleware.VerifyTokenFunc
		authToken      string
	}{
		{
			name: "success",
			mockSetup: func(m *MockUploadServiceClient) {
				m.On("CreatePresignedUrl", mock.Anything, &uploadpb.CreatePresignedUrlRequest{}).Return(&uploadpb.CreatePresignedUrlResponse{
					Message: constant.MessageOK,
					Data: &uploadpb.CreatePresignedUrlResponseData{
						VideoId:      "1",
						ThumbnailId:  "1",
						VideoUrl:     "example.com/video",
						ThumbnailUrl: "example.com/thumbnail",
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data": gin.H{
					"video_id":      "1",
					"thumbnail_id":  "1",
					"video_url":     "example.com/video",
					"thumbnail_url": "example.com/thumbnail",
				},
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
		{
			name: "grpc error",
			mockSetup: func(m *MockUploadServiceClient) {
				m.On("CreatePresignedUrl", mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
		{
			name:           "verify token grpc error",
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenError,
			authToken:      "token",
		},
		{
			name:           "invalid auth token",
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: gin.H{
				"message": constant.MessageUnauthorized,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "invalid-token",
		},
		{
			name:           "missing authorization header",
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: gin.H{
				"message": constant.MessageUnauthorized,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenSuccess, // won't be called
			authToken:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockUploadServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewUploadHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/videos/upload/presigned-url", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.authToken != "" {
				c.Request.Header.Set("Authorization", tt.authToken)
			}

			middleware.VerifyTokenMiddleware(tt.mockVerifyFunc)(c)

			if !c.IsAborted() {
				h.CreatePresignedUrl(c)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUploadHandler_UploadedWebhook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	makeReqBody := func(t *testing.T, body any) *bytes.Reader {
		t.Helper()
		switch b := body.(type) {
		case string:
			return bytes.NewReader([]byte(b))
		default:
			data, err := json.Marshal(b)
			require.NoError(t, err)
			return bytes.NewReader(data)
		}
	}

	reqBody := types.UploadedWebhookInput{
		VideoId:     "1",
		ThumbnailId: "1",
		Title:       "title",
		Description: "description",
	}

	tests := []struct {
		name           string
		body           any
		mockSetup      func(m *MockUploadServiceClient)
		expectedStatus int
		expectedJSON   gin.H
		mockVerifyFunc middleware.VerifyTokenFunc
		authToken      string
	}{
		{
			name: "success",
			body: reqBody,
			mockSetup: func(m *MockUploadServiceClient) {
				m.On("UploadedWebhook", mock.Anything, &uploadpb.UploadedWebhookRequest{
					VideoId:     reqBody.VideoId,
					ThumbnailId: reqBody.ThumbnailId,
					Title:       reqBody.Title,
					Description: reqBody.Description,
				}).Return(&uploadpb.UploadedWebhookResponse{
					Message: constant.MessageOK,
					Data:    &uploadpb.UploadedWebhookResponseData{},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data":    gin.H{},
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
		{
			name: "grpc error",
			body: reqBody,
			mockSetup: func(m *MockUploadServiceClient) {
				m.On("UploadedWebhook", mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
		{
			name:           "verify token grpc error",
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenError,
			authToken:      "token",
		},
		{
			name:           "invalid auth token",
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: gin.H{
				"message": constant.MessageUnauthorized,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "invalid-token",
		},
		{
			name:           "missing authorization header",
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: gin.H{
				"message": constant.MessageUnauthorized,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenSuccess, // won't be called
			authToken:      "",
		},
		{
			name: "invalid json",
			body: `{"title":123}`,
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data": gin.H{
					"errors": gin.H{},
				},
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
		{
			name: "missing/empty fields",
			body: types.UploadedWebhookInput{VideoId: "", ThumbnailId: "", Title: ""},
			mockSetup: func(m *MockUploadServiceClient) {
				// no call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data": gin.H{
					"errors": gin.H{
						"video_id":     []string{"video_id is required"},
						"thumbnail_id": []string{"thumbnail_id is required"},
						"title":        []string{"title is required"},
						"description":  []string{"description is required"},
					},
				},
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockUploadServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewUploadHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/videos/upload/webhook", makeReqBody(t, tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.authToken != "" {
				c.Request.Header.Set("Authorization", tt.authToken)
			}

			middleware.VerifyTokenMiddleware(tt.mockVerifyFunc)(c)

			if !c.IsAborted() {
				h.UploadedWebhook(c)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}
