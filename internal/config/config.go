package config

import (
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gofor-little/env"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
)

var conf *Config

type Config struct {
	HTTPServer *httpServer
	GRPCClient *grpcClient
}

type httpServer struct {
	Host string
	Port int
}

type grpcClient struct {
	AuthenticationServiceurl string
	UploadServiceurl         string
	Timeout                  time.Duration
}

func Init() {
	envPath := path.Join(helper.RootDir(), "..", ".env")

	if err := env.Load(envPath); err != nil {
		log.Fatal("Failed to load .env %q: %v", envPath, err)
	}

	log.Info("Loaded %q", envPath)

	port, err := strconv.Atoi(Getenv("HTTP_PORT", "4000"))

	if err != nil {
		log.Error("Invalid HTTP_PORT value %v", err)
	}

	timeout, err := strconv.Atoi(Getenv("GRPC_CLIENT_TIMEOUT_SECONDS", "5"))

	if err != nil {
		log.Error("Invalid GRPC_CLIENT_TIMEOUT_SECONDS value %v", err)
	}

	conf = &Config{
		HTTPServer: &httpServer{
			Host: Getenv("HTTP_HOST", "localhost"),
			Port: port,
		},
		GRPCClient: &grpcClient{
			AuthenticationServiceurl: Getenv("GRPC_AUTHENTICATION_SERVICE_URL", "localhost:5001"),
			UploadServiceurl:         Getenv("GRPC_UPLOAD_SERVICE_URL", "localhost:5002"),
			Timeout:                  time.Duration(timeout) * time.Second,
		},
	}
}

func GethttpServer() *httpServer {
	return conf.HTTPServer
}

func GetgrpcClient() *grpcClient {
	return conf.GRPCClient
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}
