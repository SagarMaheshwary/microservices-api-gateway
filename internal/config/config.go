package config

import (
	"os"
	"path"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
)

type Config struct {
	HTTPServer               *HTTPServer
	App                      *App
	GRPCAuthenticationClient *GRPCAuthenticationClient
	GRPCUploadClient         *GRPCUploadClient
	GRPCVideoCatalogClient   *GRPCVideoCatalogClient
	Jaeger                   *Jaeger
}

type HTTPServer struct {
	Host string
	Port int
}

type App struct {
	Env string
}

type GRPCAuthenticationClient struct {
	URL     string
	Timeout time.Duration
}

type GRPCUploadClient struct {
	URL     string
	Timeout time.Duration
}

type GRPCVideoCatalogClient struct {
	URL     string
	Timeout time.Duration
}

type Jaeger struct {
	URL string
}

type LoaderOptions struct {
	EnvPath     string
	EnvLoader   func(string) error
	FileChecker func(string) bool
}

func NewConfigWithOptions(opts LoaderOptions) *Config {
	envLoader := opts.EnvLoader
	if envLoader == nil {
		envLoader = func(path string) error { return env.Load(path) }
	}
	fileChecker := opts.FileChecker
	if fileChecker == nil {
		fileChecker = func(path string) bool {
			_, err := os.Stat(path)
			return err == nil
		}
	}

	if opts.EnvPath != "" && fileChecker(opts.EnvPath) {
		if err := envLoader(opts.EnvPath); err != nil {
			logger.Panic("Failed to load .env %q: %v", opts.EnvPath, err)
		}
		logger.Info("Loaded environment variables from %q", opts.EnvPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	return &Config{
		HTTPServer: &HTTPServer{
			Host: helper.GetEnv("HTTP_HOST", "localhost"),
			Port: helper.GetEnvInt("HTTP_PORT", 4000),
		},
		App: &App{
			Env: helper.GetEnv("APP_ENV", "development"),
		},
		GRPCAuthenticationClient: &GRPCAuthenticationClient{
			URL:     helper.GetEnv("GRPC_AUTHENTICATION_SERVICE_URL", "authentication-service:5001"),
			Timeout: helper.GetEnvDurationSeconds("GRPC_AUTHENTICATION_SERVICE_TIMEOUT_SECONDS", 3),
		},
		GRPCUploadClient: &GRPCUploadClient{
			URL:     helper.GetEnv("GRPC_UPLOAD_SERVICE_URL", "upload-service:5002"),
			Timeout: helper.GetEnvDurationSeconds("GRPC_UPLOAD_SERVICE_TIMEOUT_SECONDS", 3),
		},
		GRPCVideoCatalogClient: &GRPCVideoCatalogClient{
			URL:     helper.GetEnv("GRPC_VIDEO_CATALOG_SERVICE_URL", "video-catalog-service:5001"),
			Timeout: helper.GetEnvDurationSeconds("GRPC_VIDEO_CATALOG_SERVICE_TIMEOUT_SECONDS", 3),
		},
		Jaeger: &Jaeger{
			URL: helper.GetEnv("JAEGER_URL", "jaeger:4318"),
		},
	}
}

func NewConfig() *Config {
	return NewConfigWithOptions(LoaderOptions{
		EnvPath: path.Join(helper.GetRootDir(), "..", ".env"),
	})
}
