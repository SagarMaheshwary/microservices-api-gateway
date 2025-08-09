package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
)

type Config struct {
	HTTPServer *HTTPServer
	App        *App
	GRPCClient *GRPCClient
	Jaeger     *Jaeger
}

type HTTPServer struct {
	Host string
	Port int
}

type App struct {
	Env string
}

type GRPCClient struct {
	AuthenticationServiceURL string
	UploadServiceURL         string
	VideoCatalogServiceURL   string
	Timeout                  time.Duration
}

type Jaeger struct {
	URL string
}

func NewConfig() *Config {
	envPath := path.Join(helper.GetRootDir(), "..", ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := env.Load(envPath); err != nil {
			logger.Fatal("Failed to load .env %q: %v", envPath, err)
		}

		logger.Info("Loaded environment variables from %q", envPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	return &Config{
		HTTPServer: &HTTPServer{
			Host: getEnv("HTTP_HOST", "localhost"),
			Port: getEnvInt("HTTP_PORT", 4000),
		},
		App: &App{
			Env: getEnv("APP_ENV", "development"),
		},
		GRPCClient: &GRPCClient{
			AuthenticationServiceURL: getEnv("GRPC_AUTHENTICATION_SERVICE_URL", "authentication-service:5001"),
			UploadServiceURL:         getEnv("GRPC_UPLOAD_SERVICE_URL", "upload-service:5002"),
			VideoCatalogServiceURL:   getEnv("GRPC_VIDEO_CATALOG_SERVICE_URL", "video-catalog-service:5002"),
			Timeout:                  getEnvDurationSeconds("GRPC_CLIENT_TIMEOUT_SECONDS", 5),
		},
		Jaeger: &Jaeger{
			URL: getEnv("JAEGER_URL", "jaeger:4318"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return val
	}

	return defaultVal
}

func getEnvDurationSeconds(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
