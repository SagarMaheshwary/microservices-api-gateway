package authentication

import (
	"context"
	"errors"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

type AuthenticationService interface {
	Register(ctx context.Context, in *authpb.RegisterRequest) (*authpb.RegisterResponse, error)
	Login(ctx context.Context, in *authpb.LoginRequest) (*authpb.LoginResponse, error)
	VerifyToken(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error)
	Logout(ctx context.Context, in *authpb.LogoutRequest, token string) (*authpb.LogoutResponse, error)
	Health(ctx context.Context) error
}

type AuthenticationClient struct {
	config *config.GRPCClient
	client authpb.AuthenticationServiceClient
	health healthpb.HealthClient
}

func NewAuthenticationClient(c authpb.AuthenticationServiceClient, h healthpb.HealthClient, cfg *config.GRPCClient) *AuthenticationClient {
	return &AuthenticationClient{
		client: c,
		health: h,
		config: cfg,
	}
}

func (a *AuthenticationClient) Register(ctx context.Context, in *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	response, err := a.client.Register(ctx, in)
	if err != nil {
		logger.Error("gRPC authenticationClient.Register failed %v", err)
		return nil, err
	}

	return response, nil
}

func (a *AuthenticationClient) Login(ctx context.Context, in *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	response, err := a.client.Login(ctx, in)
	if err != nil {
		logger.Error("gRPC authenticationClient.Login failed %v", err)
		return nil, err
	}

	return response, nil
}

func (a *AuthenticationClient) VerifyToken(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	md := metadata.Pairs(constant.GRPCHeaderAuthorization, token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := a.client.VerifyToken(ctx, in)
	if err != nil {
		logger.Error("gRPC authenticationClient.VerifyToken failed %v", err)
		return nil, err
	}

	return response, nil
}

func (a *AuthenticationClient) Logout(ctx context.Context, in *authpb.LogoutRequest, token string) (*authpb.LogoutResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	md := metadata.Pairs(constant.GRPCHeaderAuthorization, token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	response, err := a.client.Logout(ctx, in)
	if err != nil {
		logger.Error("gRPC authenticationClient.Logout failed %v", err)
		return nil, err
	}

	return response, nil
}

func (a *AuthenticationClient) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, a.config.Timeout)
	defer cancel()

	res, err := a.health.Check(ctx, &healthpb.HealthCheckRequest{})
	if err != nil {
		logger.Error("Authentication gRPC health check failed! %v", err)
		return err
	}

	if res.Status == healthpb.HealthCheckResponse_NOT_SERVING {
		logger.Error("Authentication gRPC health check failed")
		return errors.New("authentication grpc health check failed")
	}

	return nil
}
