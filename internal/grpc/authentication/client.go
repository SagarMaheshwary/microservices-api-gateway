package authentication

import (
	"github.com/sagarmaheshwary/microservices-api-gateway/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	apb "github.com/sagarmaheshwary/microservices-api-gateway/proto/authentication/authentication"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	address := config.GetgrpcClient().AuthenticationServiceurl

	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		log.Error("gRPC client failed to connect on %q: %v", address, err)
	}

	log.Info("gRPC client connected on %q", address)

	Auth = &authenticationClient{
		client: apb.NewAuthenticationServiceClient(conn),
	}
}
