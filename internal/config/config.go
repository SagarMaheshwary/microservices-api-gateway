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

var Conf *Config

type Config struct {
	HTTPServer *HTTPServer
	GRPCClient *GRPCClient
}

type HTTPServer struct {
	Host string
	Port int
}

type GRPCClient struct {
	AuthenticationServiceurl string
	UploadServiceurl         string
	VideoCatalogServiceurl   string
	Timeout                  time.Duration
}

func Init() {
	envPath := path.Join(helper.GetRootDir(), "..", ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := env.Load(envPath); err != nil {
			logger.Fatal("Failed to load .env %q: %v", envPath, err)
		}

		logger.Info("Loaded environment variables from %q", envPath)
	} else {
		logger.Info(".env file not found, using system environment variables")
	}

	Conf = &Config{
		HTTPServer: &HTTPServer{
			Host: getEnv("HTTP_HOST", "localhost"),
			Port: getEnvInt("HTTP_PORT", 4000),
		},
		GRPCClient: &GRPCClient{
			AuthenticationServiceurl: getEnv("GRPC_AUTHENTICATION_SERVICE_URL", "localhost:5001"),
			UploadServiceurl:         getEnv("GRPC_UPLOAD_SERVICE_URL", "localhost:5002"),
			VideoCatalogServiceurl:   getEnv("GRPC_VIDEO_CATALOG_SERVICE_URL", "localhost:5002"),
			Timeout:                  getEnvDuration("GRPC_CLIENT_TIMEOUT_SECONDS", 5),
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

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return time.Duration(val) * time.Second
	}

	return defaultVal * time.Second
}
