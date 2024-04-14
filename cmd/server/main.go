package main

import (
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
)

func main() {
	log.Init()
	config.Init()

	authrpc.Connect()
	http.Connect()
}
