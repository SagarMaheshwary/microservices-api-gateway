package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewServer(cfg *config.HTTPServer, env string, grpcClients types.GRPCClients) *http.Server {
	address := fmt.Sprintf("%v:%d", cfg.Host, cfg.Port)

	if env == "development" {
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

	authHandler := handler.NewAuthHandler(grpcClients.AuthClient)
	healthHandler := handler.NewHealthHandler(grpcClients)
	uploadHandler := handler.NewUploadHandler(grpcClients.UploadClient)
	videoCatalogHandler := handler.NewVideoCatalogHandler(grpcClients.VideoCatalogClient)

	r.GET("/health", healthHandler.CheckAll)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		authenticated := auth.Group("/", middleware.VerifyTokenMiddleware(grpcClients.AuthClient))
		{
			authenticated.GET("/profile", authHandler.Profile)
			authenticated.POST("/logout", authHandler.Logout)
		}
	}

	videos := r.Group("/videos")
	{
		videos.GET("", videoCatalogHandler.FindAll)
		videos.GET("/:id", videoCatalogHandler.FindById)

		authenticated := videos.Group("/", middleware.VerifyTokenMiddleware(grpcClients.AuthClient))
		{
			authenticated.POST("/upload/presigned-url", uploadHandler.CreatePresignedUrl)
			authenticated.POST("/upload/webhook", uploadHandler.UploadedWebhook)
		}
	}

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
