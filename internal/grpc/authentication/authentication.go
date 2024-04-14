package authentication

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	cons "github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"google.golang.org/grpc/metadata"
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

func (a *authenticationClient) VerifyToken(data *pb.VerifyTokenRequest, token string) (*pb.VerifyTokenResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	md := metadata.Pairs(cons.HeaderAuthorization, token)
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	response, err := a.client.VerifyToken(ctx, data)

	if err != nil {
		log.Error("gRPC authenticationClient.VerifyToken failed %v", err)

		return nil, err
	}

	return response, nil
}

func (a *authenticationClient) Logout(data *pb.LogoutRequest, token string) (*pb.LogoutResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	md := metadata.Pairs(cons.HeaderAuthorization, token)
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	response, err := a.client.Logout(ctx, data)

	if err != nil {
		log.Error("gRPC authenticationClient.Logout failed %v", err)

		return nil, err
	}

	return response, nil
}
