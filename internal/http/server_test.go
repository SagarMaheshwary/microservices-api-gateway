package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	myhttp "github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthHandler struct {
	mock.Mock
}

func (m *mockAuthHandler) Register(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "register")
}
func (m *mockAuthHandler) Login(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "login")
}
func (m *mockAuthHandler) Profile(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "profile")
}
func (m *mockAuthHandler) Logout(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "logout")
}

type mockHealthHandler struct{ mock.Mock }

func (m *mockHealthHandler) CheckAll(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "health")
}

type mockUploadHandler struct{ mock.Mock }

func (m *mockUploadHandler) CreatePresignedUrl(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "presigned")
}
func (m *mockUploadHandler) UploadedWebhook(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "webhook")
}

type mockVideoCatalogHandler struct{ mock.Mock }

func (m *mockVideoCatalogHandler) FindAll(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "videos")
}
func (m *mockVideoCatalogHandler) FindById(c *gin.Context) {
	m.Called(c)
	c.String(http.StatusOK, "video by id")
}

func TestNewRouter_AllRoutes_WithTestifyMocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authMock := new(mockAuthHandler)
	healthMock := new(mockHealthHandler)
	uploadMock := new(mockUploadHandler)
	videoMock := new(mockVideoCatalogHandler)

	verifyToken := func(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
		return &authpb.VerifyTokenResponse{
			Message: constant.MessageOK,
			Data: &authpb.VerifyTokenResponseData{
				User: &authpb.User{
					Id:    1,
					Name:  "name",
					Email: "name@gmail.com",
				},
			},
		}, nil
	}

	router := myhttp.NewRouter(myhttp.RouterConfig{
		Env:                 "test",
		AuthHandler:         authMock,
		HealthHandler:       healthMock,
		UploadHandler:       uploadMock,
		VideoCatalogHandler: videoMock,
		VerifyToken:         verifyToken,
	})

	tests := []struct {
		name         string
		method       string
		target       string
		body         string
		wantStatus   int
		wantBody     string
		requiresAuth bool
		mockCalls    func()
	}{
		{"health", "GET", "/health", "", 200, "health", false, func() { healthMock.On("CheckAll", mock.Anything).Once() }},
		{"metrics", "GET", "/metrics", "", 200, "", false, func() {}},

		{"register", "POST", "/auth/register", "", 200, "register", false, func() { authMock.On("Register", mock.Anything).Once() }},
		{"login", "POST", "/auth/login", "", 200, "login", true, func() { authMock.On("Login", mock.Anything).Once() }},
		{"profile", "GET", "/auth/profile", "", 200, "profile", true, func() { authMock.On("Profile", mock.Anything).Once() }},
		{"logout", "POST", "/auth/logout", "", 200, "logout", true, func() { authMock.On("Logout", mock.Anything).Once() }},

		{"videos list", "GET", "/videos", "", 200, "videos", false, func() { videoMock.On("FindAll", mock.Anything).Once() }},
		{"video by id", "GET", "/videos/123", "", 200, "video by id", false, func() { videoMock.On("FindById", mock.Anything).Once() }},
		{"create presigned", "POST", "/videos/upload/presigned-url", "", 200, "presigned", true, func() { uploadMock.On("CreatePresignedUrl", mock.Anything).Once() }},
		{"upload webhook", "POST", "/videos/upload/webhook", "", 200, "webhook", true, func() { uploadMock.On("UploadedWebhook", mock.Anything).Once() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockCalls()

			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.target, strings.NewReader(tt.body))
			} else {
				req = httptest.NewRequest(tt.method, tt.target, nil)
			}

			if tt.requiresAuth {
				req.Header.Set("Authorization", "Bearer faketoken")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.Contains(t, w.Body.String(), tt.wantBody)
			}

			mock.AssertExpectationsForObjects(t, authMock, healthMock, uploadMock, videoMock)
		})
	}
}
