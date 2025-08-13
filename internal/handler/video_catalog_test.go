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
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var dummyVideo = &videocatalogpb.Video{
	Id:           1,
	Title:        "title",
	Description:  "description",
	ThumbnailUrl: "example.com/thumbnail",
	PublishedAt:  now,
	Duration:     120,
	Resolution:   "1920x1080",
	User: &videocatalogpb.User{
		Id:    1,
		Name:  "name",
		Image: nil,
	},
}

func TestVideoCatalogHandler_FindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(m *MockVideoCatalogServiceClient)
		expectedStatus int
		expectedJSON   gin.H
	}{
		{
			name: "success",
			mockSetup: func(m *MockVideoCatalogServiceClient) {
				m.On("FindAll", mock.Anything, &videocatalogpb.FindAllRequest{}).
					Return(&videocatalogpb.FindAllResponse{
						Message: constant.MessageOK,
						Data: &videocatalogpb.FindAllResponseData{
							Videos: []*videocatalogpb.Video{dummyVideo},
						},
					}, nil).
					Once()
			},
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data": gin.H{
					"videos": []any{
						gin.H{
							"id":            1,
							"title":         "title",
							"description":   "description",
							"thumbnail_url": "example.com/thumbnail",
							"published_at":  now,
							"duration":      120,
							"resolution":    "1920x1080",
							"user": gin.H{
								"id":   1,
								"name": "name",
							},
						},
					},
				},
			},
		},
		{
			name: "grpc error",
			mockSetup: func(m *MockVideoCatalogServiceClient) {
				m.On("FindAll", mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).
					Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockVideoCatalogServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewVideoCatalogHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/videos", nil)
			c.Request.Header.Set("Content-Type", "application/json")

			h.FindAll(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestVideoCatalogHandler_FindById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockSetup      func(m *MockVideoCatalogServiceClient)
		videoId        string
		expectedStatus int
		expectedJSON   gin.H
	}{
		{
			name: "success",
			mockSetup: func(m *MockVideoCatalogServiceClient) {
				m.On("FindById", mock.Anything, &videocatalogpb.FindByIdRequest{
					Id: 1,
				}).
					Return(&videocatalogpb.FindByIdResponse{
						Message: constant.MessageOK,
						Data: &videocatalogpb.FindByIdResponseData{
							Video: dummyVideo,
						},
					}, nil).
					Once()
			},
			videoId:        "1",
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data": gin.H{
					"video": gin.H{
						"id":            1,
						"title":         "title",
						"description":   "description",
						"thumbnail_url": "example.com/thumbnail",
						"published_at":  now,
						"duration":      120,
						"resolution":    "1920x1080",
						"user": gin.H{
							"id":   1,
							"name": "name",
						},
					},
				},
			},
		},
		{
			name: "grpc error",
			mockSetup: func(m *MockVideoCatalogServiceClient) {
				m.On("FindById", mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).
					Once()
			},
			videoId:        "1",
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
		},
		{
			name: "invalid or missing id",
			mockSetup: func(m *MockVideoCatalogServiceClient) {
				// No gRPC call expected for invalid ID
			},
			videoId:        "", // Simulating missing param
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data":    gin.H{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockVideoCatalogServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewVideoCatalogHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest(http.MethodGet, "/videos/"+tt.videoId, nil)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			if tt.videoId != "" {
				c.Params = gin.Params{
					gin.Param{Key: "id", Value: tt.videoId},
				}
			}

			h.FindById(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}
