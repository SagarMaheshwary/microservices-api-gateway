package authentication

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-api-gateway/proto/authentication/authentication"
)

var Auth *authenticationClient

type authenticationClient struct {
	client pb.AuthenticationServiceClient
}

func (a *authenticationClient) Register(data *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := a.client.Register(ctx, data)

	if err != nil {
		log.Error("gRPC authenticationClient.Register failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) Login(data *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := a.client.Login(ctx, data)

	if err != nil {
		log.Error("gRPC authenticationClient.Login failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) VerifyToken(data *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := a.client.VerifyToken(ctx, data)

	if err != nil {
		log.Error("gRPC authenticationClient.VerifyToken failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) Logout(data *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := a.client.Logout(ctx, data)

	if err != nil {
		log.Error("gRPC authenticationClient.Logout failed %v", err)

		return nil, err
	}

	return response, nil
}
