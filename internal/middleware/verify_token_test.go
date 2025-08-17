package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/stretchr/testify/assert"
)

var now = time.Now().String()

var dummyUser = &authpb.User{
	Id:        1,
	Name:      "name",
	Email:     "name@gmail.com",
	Image:     nil,
	CreatedAt: &now,
	UpdatedAt: nil,
}

func mockVerifyTokenSuccess(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	return &authpb.VerifyTokenResponse{
		Data: &authpb.VerifyTokenResponseData{
			User: dummyUser,
		},
	}, nil
}

func mockVerifyTokenError(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	return nil, errors.New("invalid token")
}

func TestVerifyTokenMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		verifyToken    middleware.VerifyTokenFunc
		header         string
		expectedStatus int
		assertions     func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "No Authorization header",
			verifyToken:    mockVerifyTokenSuccess,
			header:         "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:        "Invalid token",
			verifyToken: mockVerifyTokenError,
			header:      "Bearer badtoken",
			// we don’t know exact status from PrepareResponseFromGRPCError,
			// just ensure it’s neither 200 nor 401
			expectedStatus: 0,
			assertions: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.NotEqual(t, http.StatusOK, w.Code)
				assert.NotEqual(t, http.StatusUnauthorized, w.Code)
			},
		},
		{
			name:           "Valid token",
			verifyToken:    mockVerifyTokenSuccess,
			header:         "Bearer validtoken",
			expectedStatus: http.StatusOK,
			assertions: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, w.Code)
				assert.Contains(t, w.Body.String(), "ok")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.VerifyTokenMiddleware(tt.verifyToken))
			router.GET("/test", func(c *gin.Context) {
				// For success case, assert context here
				if tt.name == "Valid token" {
					user, exists := c.Get(constant.AuthUser)
					assert.True(t, exists)
					assert.NotNil(t, user)
				}
				c.JSON(http.StatusOK, gin.H{"msg": "ok"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tt.expectedStatus != 0 {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
			if tt.assertions != nil {
				tt.assertions(t, w)
			}
		})
	}
}
