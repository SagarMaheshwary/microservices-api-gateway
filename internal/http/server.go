package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/router"
)

func Connect() {
	c := config.Conf.HTTPServer
	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	r := gin.Default()

	r.Use(middleware.PrometheusMiddleware())
	r.Use(middleware.CORSMiddleware()) //@TODO: should it be in code or reverse proxy?

	router.InitRoutes(r)

	s := &http.Server{
		Addr:    address,
		Handler: r,
	}

	logger.Info("HTTP server started on %v", address)

	if err := s.ListenAndServe(); err != nil {
		logger.Error("HTTP server failed to start %v", err)
	}
}
