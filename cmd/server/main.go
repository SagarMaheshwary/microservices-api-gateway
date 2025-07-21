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
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	httpserver "github.com/sagarmaheshwary/microservices-api-gateway/internal/http"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/jaeger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
	"google.golang.org/grpc"
)

func main() {
	logger.Init()
	config.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	shutdownJaeger := jaeger.Init(ctx)
	prometheus.RegisterMetrics()

	authConn := mustInitClient("Authentication", authrpc.InitClient)
	defer authConn.Close()

	uploadConn := mustInitClient("Upload", uploadrpc.InitClient)
	defer uploadConn.Close()

	videocatalogConn := mustInitClient("VideoCatalog", videocatalogrpc.InitClient)
	defer videocatalogConn.Close()

	httpServer := httpserver.NewServer()
	go func() {
		if err := httpserver.Serve(httpServer); err != nil && err != http.ErrServerClosed {
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

func mustInitClient(name string, initFunc func(ctx context.Context) (*grpc.ClientConn, error)) *grpc.ClientConn {
	conn, err := initFunc(context.Background())
	if err != nil {
		logger.Error("Failed to init %s client: %v", name, err)
		os.Exit(constant.ExitFailure)
	}

	return conn
}
