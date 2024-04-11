package main

import (
	"github.com/sagarmaheshwary/microservices-api-gateway/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
)

func main() {
	log.Init()
	config.Init()

	authentication.Connect()
	http.Connect()
}
