package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	auth "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	upload "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalog "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video-catalog"
	httpserver "github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

func main() {
	logger.Init()
	cfg := config.NewConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdownJaeger := jaeger.Init(ctx, cfg.Jaeger.URL)
	prometheus.RegisterMetrics()

	authClient, authConn, err := auth.NewClient(ctx, &auth.InitClientOptions{Config: cfg.GRPCAuthenticationClient})
	if err != nil {
		logger.Error("Failed to connect to auth client: %v", err)
		os.Exit(constant.ExitFailure)
	}
	defer authConn.Close()

	uploadClient, uploadConn, err := upload.NewClient(ctx, &upload.InitClientOptions{Config: cfg.GRPCUploadClient})
	if err != nil {
		logger.Error("Failed to connect to upload client: %v", err)
		os.Exit(constant.ExitFailure)
	}
	defer uploadConn.Close()

	videoCatalogClient, videoCatalogConn, err := videocatalog.NewClient(ctx, &videocatalog.InitClientOptions{Config: cfg.GRPCVideoCatalogClient})
	if err != nil {
		logger.Error("Failed to connect to video catalog client: %v", err)
		os.Exit(constant.ExitFailure)
	}
	defer videoCatalogConn.Close()

	httpServer := httpserver.NewServer(
		cfg.HTTPServer,
		cfg.App.Env,
		types.GRPCClients{
			AuthClient:         authClient,
			UploadClient:       uploadClient,
			VideoCatalogClient: videoCatalogClient,
		},
	)

	go func() {
		if err := httpserver.Serve(httpServer, httpServer.ListenAndServe); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	logger.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdownJaeger(shutdownCtx); err != nil {
		logger.Warn("failed to shutdown jaeger tracer: %v", err)
	}

	shutdownCtx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Warn("Http server shutdown error: %v", err)
	}

	logger.Info("Shutdown complete")
}
