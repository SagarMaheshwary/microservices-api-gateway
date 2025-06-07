package main

import (
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	userrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
)

func main() {
	logger.Init()
	config.Init()

	authrpc.Connect()
	userrpc.Connect()
	videocatalogrpc.Connect()

	prometheus.Init()
	jaeger.Init()

	http.Connect()
}
