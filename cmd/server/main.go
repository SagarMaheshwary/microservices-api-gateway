package main

import (
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	arpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	urpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	vcrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
)

func main() {
	log.Init()
	config.Init()

	arpc.Connect()
	urpc.Connect()
	vcrpc.Connect()

	http.Connect()
}
