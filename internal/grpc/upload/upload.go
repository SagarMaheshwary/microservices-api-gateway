package upload

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	cons "github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	pb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"google.golang.org/grpc/metadata"
)

var Upload *uploadClient

type uploadClient struct {
	client pb.UploadServiceClient
}

func (u *uploadClient) CreatePresignedUrl(data *pb.CreatePresignedUrlRequest) (*pb.CreatePresignedUrlResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := u.client.CreatePresignedUrl(ctx, data)

	if err != nil {
		log.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)

		return nil, err
	}

	return response, nil
}

func (u *uploadClient) UploadedWebhook(data *pb.UploadedWebhookRequest, userId string) (*pb.UploadedWebhookResponse, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	md := metadata.Pairs(cons.HeaderUserId, userId)
	ctx := metadata.NewOutgoingContext(ctxTimeout, md)

	response, err := u.client.UploadedWebhook(ctx, data)

	if err != nil {
		log.Error("gRPC uploadClient.CreatePresignedUrl failed %v", err)

		return nil, err
	}

	return response, nil
}
