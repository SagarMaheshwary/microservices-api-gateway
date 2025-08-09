package types

import (
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video-catalog"
)

type GRPCClients struct {
	AuthClient         authrpc.AuthenticationService
	UploadClient       uploadrpc.UploadService
	VideoCatalogClient videocatalogrpc.VideoCatalogService
}
