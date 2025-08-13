package authentication_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	auth "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
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

func TestAuthenticationClient_Register(t *testing.T) {
	req := &authpb.RegisterRequest{
		Email:    "name@gmail.com",
		Password: "password",
	}
	res := &authpb.RegisterResponse{
		Message: constant.MessageOK,
		Data:    &authpb.RegisterResponseData{Token: "token", User: dummyUser},
	}

	cfg := &config.GRPCAuthenticationClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *authpb.RegisterResponse
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
			name:       "generic failure",
			mockReturn: nil,
			mockErr:    errors.New("grpc failure"),
			expectErr:  true,
		},
		{
			name:       "email taken",
			mockReturn: nil,
			mockErr:    status.Error(codes.AlreadyExists, "email taken"),
			expectErr:  true,
			expectGRPC: codes.AlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAuthenticationServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("Register", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := auth.NewAuthenticationClient(mockClient, mockHealth, cfg)

			got, err := c.Register(context.Background(), req)

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

func TestAuthenticationClient_Login(t *testing.T) {
	req := &authpb.LoginRequest{
		Email:    "name@gmail.com",
		Password: "password",
	}
	res := &authpb.LoginResponse{
		Message: constant.MessageOK,
		Data:    &authpb.LoginResponseData{Token: "token", User: dummyUser},
	}

	cfg := &config.GRPCAuthenticationClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *authpb.LoginResponse
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
			name:       "unauthenticated (wrong credentials)",
			mockReturn: nil,
			mockErr:    status.Error(codes.Unauthenticated, "invalid credentials"),
			expectErr:  true,
			expectGRPC: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAuthenticationServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("Login", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := auth.NewAuthenticationClient(mockClient, mockHealth, cfg)

			got, err := c.Login(context.Background(), req)

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

func TestAuthenticationClient_VerifyToken(t *testing.T) {
	req := &authpb.VerifyTokenRequest{}
	res := &authpb.VerifyTokenResponse{
		Message: constant.MessageOK,
		Data:    &authpb.VerifyTokenResponseData{User: dummyUser},
	}

	cfg := &config.GRPCAuthenticationClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *authpb.VerifyTokenResponse
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
			name:       "unauthenticated (wrong/expired token)",
			mockReturn: nil,
			mockErr:    status.Error(codes.Unauthenticated, "invalid credentials"),
			expectErr:  true,
			expectGRPC: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAuthenticationServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("VerifyToken", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := auth.NewAuthenticationClient(mockClient, mockHealth, cfg)

			got, err := c.VerifyToken(context.Background(), req, "token")

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

func TestAuthenticationClient_Logout(t *testing.T) {
	req := &authpb.LogoutRequest{}
	res := &authpb.LogoutResponse{
		Message: constant.MessageOK,
		Data:    &authpb.LogoutResponseData{},
	}

	cfg := &config.GRPCAuthenticationClient{Timeout: 2 * time.Second}

	tests := []struct {
		name       string
		mockReturn *authpb.LogoutResponse
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
			name:       "unauthenticated (wrong/expired token)",
			mockReturn: nil,
			mockErr:    status.Error(codes.Unauthenticated, "invalid credentials"),
			expectErr:  true,
			expectGRPC: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockAuthenticationServiceClient)
			mockHealth := new(MockHealthClient)

			mockClient.On("Logout", mock.Anything, req).
				Return(tt.mockReturn, tt.mockErr).
				Once()

			c := auth.NewAuthenticationClient(mockClient, mockHealth, cfg)

			got, err := c.Logout(context.Background(), req, "token")

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
	cfg := &config.GRPCAuthenticationClient{Timeout: 2 * time.Second}

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
			expectMsg: "authentication grpc health check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUpload := new(MockAuthenticationServiceClient)
			mockHealth := new(MockHealthClient)

			client := auth.NewAuthenticationClient(mockUpload, mockHealth, cfg)

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
