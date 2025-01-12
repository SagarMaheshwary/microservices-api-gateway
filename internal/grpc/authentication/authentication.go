package authentication

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

var Auth *authenticationClient

type authenticationClient struct {
	client authpb.AuthenticationServiceClient
	health healthpb.HealthClient
}

func (a *authenticationClient) Register(data *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := a.client.Register(ctx, data)

	if err != nil {
		logger.Error("gRPC authenticationClient.Register failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) Login(data *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	response, err := a.client.Login(ctx, data)

	if err != nil {
		logger.Error("gRPC authenticationClient.Login failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) VerifyToken(data *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	md := metadata.Pairs(constant.HeaderAuthorization, token)
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	response, err := a.client.VerifyToken(ctx, data)

	if err != nil {
		logger.Error("gRPC authenticationClient.VerifyToken failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) Logout(data *authpb.LogoutRequest, token string) (*authpb.LogoutResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.Conf.GRPCClient.Timeout)

	defer cancel()

	md := metadata.Pairs(constant.HeaderAuthorization, token)
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	response, err := a.client.Logout(ctx, data)

	if err != nil {
		logger.Error("gRPC authenticationClient.Logout failed %v", err)

		return nil, err
	}

	return response, nil
}
