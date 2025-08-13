package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
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

func TestAuthenticationHandler_Register(t *testing.T) {
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

	tests := []struct {
		name           string
		body           any
		mockSetup      func(m *MockAuthenticationServiceClient)
		expectedStatus int
		expectedJSON   gin.H
	}{
		{
			name: "success",
			body: types.RegisterInput{
				Name:     dummyUser.Name,
				Email:    dummyUser.Email,
				Password: "secret",
			},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				m.On("Register", mock.Anything, &authpb.RegisterRequest{
					Name:     dummyUser.Name,
					Email:    dummyUser.Email,
					Password: "secret",
				}).Return(&authpb.RegisterResponse{
					Message: constant.MessageOK,
					Data: &authpb.RegisterResponseData{
						Token: "token",
						User:  dummyUser,
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusCreated,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data": gin.H{
					"token": "token",
					"user": gin.H{
						"id":         1,
						"name":       "name",
						"email":      "name@gmail.com",
						"created_at": now,
					},
				},
			},
		},
		{
			name: "invalid json",
			body: `{"name":123}`,
			mockSetup: func(m *MockAuthenticationServiceClient) {
				// no call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data": gin.H{
					"errors": gin.H{},
				},
			},
		},
		{
			name: "grpc error",
			body: types.RegisterInput{
				Name:     dummyUser.Name,
				Email:    dummyUser.Email,
				Password: "secret",
			},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				m.On("Register", mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
		},
		{
			name: "empty fields",
			body: types.RegisterInput{Name: "", Email: "", Password: ""},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				// no call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data": gin.H{
					"errors": gin.H{
						"name":     []string{"name is required"},
						"email":    []string{"email is required"},
						"password": []string{"password is required"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthenticationServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewAuthHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", makeReqBody(t, tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.Register(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestAuthenticationHandler_Login(t *testing.T) {
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

	tests := []struct {
		name           string
		body           any
		expectedStatus int
		expectedJSON   gin.H
		mockSetup      func(m *MockAuthenticationServiceClient)
	}{
		{
			name: "success",
			body: types.LoginInput{Email: dummyUser.Email, Password: "secret"},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				m.On("Login", mock.Anything, &authpb.LoginRequest{
					Email:    dummyUser.Email,
					Password: "secret",
				}).Return(&authpb.LoginResponse{
					Message: constant.MessageOK,
					Data: &authpb.LoginResponseData{
						Token: "token",
						User:  dummyUser,
					},
				}, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data": gin.H{
					"token": "token",
					"user": gin.H{
						"id":         1,
						"name":       "name",
						"email":      "name@gmail.com",
						"created_at": now,
					},
				},
			},
		},
		{
			name: "invalid json",
			body: `{"email":123}`,
			mockSetup: func(m *MockAuthenticationServiceClient) {
				// no call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data": gin.H{
					"errors": gin.H{},
				},
			},
		},
		{
			name: "grpc error",
			body: types.LoginInput{Email: dummyUser.Email, Password: "secret"},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				m.On("Login", mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
		},
		{
			name: "missing/empty fields",
			body: types.LoginInput{Email: ""},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				// no call expected
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: gin.H{
				"message": constant.MessageBadRequest,
				"data": gin.H{
					"errors": gin.H{
						"email":    []string{"email is required"},
						"password": []string{"password is required"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthenticationServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewAuthHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", makeReqBody(t, tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			h.Login(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestAuthenticationHandler_Profile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedJSON   gin.H
		mockVerifyFunc middleware.VerifyTokenFunc
		authToken      string
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data": gin.H{
					"user": gin.H{
						"id":         1,
						"name":       "name",
						"email":      "name@gmail.com",
						"created_at": now,
					},
				},
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
			mockVerifyFunc: mockVerifyTokenSuccess, // won't be called
			authToken:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthenticationServiceClient)

			h := handler.NewAuthHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.authToken != "" {
				c.Request.Header.Set("Authorization", tt.authToken)
			}

			middleware.VerifyTokenMiddleware(tt.mockVerifyFunc)(c)

			if !c.IsAborted() {
				h.Profile(c)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestAuthenticationHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		expectedStatus int
		expectedJSON   gin.H
		mockSetup      func(m *MockAuthenticationServiceClient)
		mockVerifyFunc middleware.VerifyTokenFunc
		authToken      string
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
			expectedJSON: gin.H{
				"message": constant.MessageOK,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				m.On("Logout", mock.Anything, &authpb.LogoutRequest{}, "token").Return(&authpb.LogoutResponse{
					Message: constant.MessageOK,
					Data:    &authpb.LogoutResponseData{},
				}, nil).Once()
			},
			mockVerifyFunc: mockVerifyTokenSuccess,
			authToken:      "token",
		},
		{
			name:           "grpc error",
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: gin.H{
				"message": constant.MessageInternalServerError,
				"data":    gin.H{},
			},
			mockSetup: func(m *MockAuthenticationServiceClient) {
				m.On("Logout", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, errors.New("grpc failed")).Once()
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
			mockSetup: func(m *MockAuthenticationServiceClient) {
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
			mockSetup: func(m *MockAuthenticationServiceClient) {
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
			mockSetup: func(m *MockAuthenticationServiceClient) {
				// no call expected
			},
			mockVerifyFunc: mockVerifyTokenSuccess, // won't be called
			authToken:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockAuthenticationServiceClient)
			tt.mockSetup(mockSvc)

			h := handler.NewAuthHandler(mockSvc)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.authToken != "" {
				c.Request.Header.Set("Authorization", tt.authToken)
			}

			middleware.VerifyTokenMiddleware(tt.mockVerifyFunc)(c)

			if !c.IsAborted() {
				h.Logout(c)
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			expectedBody, err := json.Marshal(tt.expectedJSON)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedBody), w.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}
