package video_catalog

import (
	"context"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	vcpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
)

var VideoCatalog *videoCatalogClient

type videoCatalogClient struct {
	client vcpb.VideoCatalogServiceClient
}

func (v *videoCatalogClient) FindAll(data *vcpb.FindAllRequest) (*vcpb.FindAllResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := v.client.FindAll(ctx, data)

	if err != nil {
		log.Error("gRPC videoCatalogClient.FindAll failed %v", err)

		return nil, err
	}

	return response, nil
}

func (v *videoCatalogClient) FindById(data *vcpb.FindByIdRequest) (*vcpb.FindByIdResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.GetgrpcClient().Timeout)

	defer cancel()

	response, err := v.client.FindById(ctx, data)

	if err != nil {
		log.Error("gRPC videoCatalogClient.FindById failed %v", err)

		return nil, err
	}

	return response, nil
}
