package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/router"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewServer() *http.Server {
	c := config.Conf.HTTPServer
	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	if config.Conf.App.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(
		gin.Recovery(),
		middleware.ZerologMiddleware(),
		otelgin.Middleware(constant.ServiceName),
		middleware.PrometheusMiddleware(),
		middleware.CORSMiddleware(), //@TODO: should it be in code or reverse proxy?
	)

	router.InitRoutes(r)

	server := &http.Server{
		Addr:    address,
		Handler: r,
	}

	return server
}

func Serve(server *http.Server) error {
	logger.Info("Starting HTTP server on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("HTTP server failed to start %v", err)
	}

	return nil
}
