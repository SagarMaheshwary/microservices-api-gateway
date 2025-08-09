package authentication_test

import (
	"context"

	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type MockAuthenticationServiceClient struct {
	mock.Mock
}

func (m *MockAuthenticationServiceClient) Register(ctx context.Context, in *authpb.RegisterRequest, opts ...grpc.CallOption) (*authpb.RegisterResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.RegisterResponse), nil
}

func (m *MockAuthenticationServiceClient) Login(ctx context.Context, in *authpb.LoginRequest, opts ...grpc.CallOption) (*authpb.LoginResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.LoginResponse), nil
}

func (m *MockAuthenticationServiceClient) VerifyToken(ctx context.Context, in *authpb.VerifyTokenRequest, opts ...grpc.CallOption) (*authpb.VerifyTokenResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.VerifyTokenResponse), nil
}

func (m *MockAuthenticationServiceClient) Logout(ctx context.Context, in *authpb.LogoutRequest, opts ...grpc.CallOption) (*authpb.LogoutResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.LogoutResponse), nil
}

func (m *MockAuthenticationServiceClient) Health(ctx context.Context, opts ...grpc.CallOption) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

type MockHealthClient struct {
	mock.Mock
	healthpb.HealthClient
}

func (m *MockHealthClient) Check(ctx context.Context, in *healthpb.HealthCheckRequest, opts ...grpc.CallOption) (*healthpb.HealthCheckResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*healthpb.HealthCheckResponse), args.Error(1)
}
