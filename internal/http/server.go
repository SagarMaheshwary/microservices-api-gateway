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

type AuthHandler interface {
	Register(*gin.Context)
	Login(*gin.Context)
	Profile(*gin.Context)
	Logout(*gin.Context)
}

type HealthHandler interface {
	CheckAll(*gin.Context)
}

type UploadHandler interface {
	CreatePresignedUrl(*gin.Context)
	UploadedWebhook(*gin.Context)
}

type VideoCatalogHandler interface {
	FindAll(*gin.Context)
	FindById(*gin.Context)
}

type RouterConfig struct {
	Env                 string
	AuthHandler         AuthHandler
	HealthHandler       HealthHandler
	UploadHandler       UploadHandler
	VideoCatalogHandler VideoCatalogHandler
	VerifyToken         middleware.VerifyTokenFunc
	Middlewares         []gin.HandlerFunc
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	switch cfg.Env {
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Default middleware if not overridden
	if len(cfg.Middlewares) == 0 {
		cfg.Middlewares = []gin.HandlerFunc{
			gin.Recovery(),
			middleware.ZerologMiddleware(),
			otelgin.Middleware(constant.ServiceName),
			middleware.PrometheusMiddleware(),
			middleware.CORSMiddleware(),
		}
	}

	r.Use(cfg.Middlewares...)

	r.GET("/health", cfg.HealthHandler.CheckAll)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := r.Group("/auth")
	{
		auth.POST("/register", cfg.AuthHandler.Register)
		auth.POST("/login", cfg.AuthHandler.Login)

		authenticated := auth.Group("/", middleware.VerifyTokenMiddleware(cfg.VerifyToken))
		{
			authenticated.GET("/profile", cfg.AuthHandler.Profile)
			authenticated.POST("/logout", cfg.AuthHandler.Logout)
		}
	}

	videos := r.Group("/videos")
	{
		videos.GET("", cfg.VideoCatalogHandler.FindAll)
		videos.GET("/:id", cfg.VideoCatalogHandler.FindById)

		authenticated := videos.Group("/", middleware.VerifyTokenMiddleware(cfg.VerifyToken))
		{
			authenticated.POST("/upload/presigned-url", cfg.UploadHandler.CreatePresignedUrl)
			authenticated.POST("/upload/webhook", cfg.UploadHandler.UploadedWebhook)
		}
	}

	return r
}

func NewServer(cfg *config.HTTPServer, env string, grpcClients types.GRPCClients) *http.Server {
	router := NewRouter(RouterConfig{
		Env:                 env,
		AuthHandler:         handler.NewAuthHandler(grpcClients.AuthClient),
		HealthHandler:       handler.NewHealthHandler(grpcClients),
		UploadHandler:       handler.NewUploadHandler(grpcClients.UploadClient),
		VideoCatalogHandler: handler.NewVideoCatalogHandler(grpcClients.VideoCatalogClient),
		VerifyToken:         grpcClients.AuthClient.VerifyToken,
	})

	address := fmt.Sprintf("%v:%d", cfg.Host, cfg.Port)

	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}

func Serve(server *http.Server, listen func() error) error {
	logger.Info("Starting HTTP server on %s", server.Addr)

	if err := listen(); err != nil {
		return fmt.Errorf("HTTP server failed to start %v", err)
	}
	return nil
}
