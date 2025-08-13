package handler_test

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockAuthenticationServiceClient struct {
	mock.Mock
}

func (m *MockAuthenticationServiceClient) Register(ctx context.Context, in *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.RegisterResponse), nil
}

func (m *MockAuthenticationServiceClient) Login(ctx context.Context, in *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	args := m.Called(ctx, in)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.LoginResponse), nil
}

func (m *MockAuthenticationServiceClient) VerifyToken(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	args := m.Called(ctx, in, token)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.VerifyTokenResponse), nil
}

func (m *MockAuthenticationServiceClient) Logout(ctx context.Context, in *authpb.LogoutRequest, token string) (*authpb.LogoutResponse, error) {
	args := m.Called(ctx, in, token)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	return args.Get(0).(*authpb.LogoutResponse), nil
}

func (m *MockAuthenticationServiceClient) Health(ctx context.Context) error {
	args := m.Called(ctx)

	if err := args.Error(0); err != nil {
		return err
	}

	return nil
}

func mockVerifyTokenSuccess(_ context.Context, _ *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	if token != "token" {
		return nil, status.Errorf(codes.Unauthenticated, constant.MessageUnauthorized)
	}
	return &authpb.VerifyTokenResponse{
		Message: constant.MessageOK,
		Data: &authpb.VerifyTokenResponseData{
			User: dummyUser,
		},
	}, nil
}

func mockVerifyTokenError(_ context.Context, _ *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	return nil, errors.New("grpc failed")
}
