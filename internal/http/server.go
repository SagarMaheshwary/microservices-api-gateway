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

func Connect() {
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

	s := &http.Server{
		Addr:    address,
		Handler: r,
	}

	logger.Info("HTTP server started on %v", address)

	if err := s.ListenAndServe(); err != nil {
		logger.Error("HTTP server failed to start %v", err)
	}
}
