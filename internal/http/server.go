package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/router"
)

func Connect() {
	c := config.GethttpServer()
	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	r := gin.Default()

	router.InitRoutes(r)

	s := &http.Server{
		Addr:    address,
		Handler: r,
	}

	log.Info("HTTP server started on %v", address)

	if err := s.ListenAndServe(); err != nil {
		log.Error("HTTP server failed to start %v", err)
	}
}
