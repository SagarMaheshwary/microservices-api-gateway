package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video-catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

type HealthHandler struct {
	authClient         authrpc.AuthenticationService
	uploadClient       uploadrpc.UploadService
	videoCatalogClient videocatalogrpc.VideoCatalogService
}

func NewHealthHandler(clients types.GRPCClients) *HealthHandler {
	return &HealthHandler{
		authClient:         clients.AuthClient,
		uploadClient:       clients.UploadClient,
		videoCatalogClient: clients.VideoCatalogClient,
	}
}

func (h *HealthHandler) CheckAll(c *gin.Context) {
	if !getServicesHealthStatus(c, h) {
		prometheus.ServiceHealth.Set(0)
		response := helper.PrepareResponse("Some services are not available!", gin.H{
			"status": "degraded",
		})
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	prometheus.ServiceHealth.Set(1)
	response := helper.PrepareResponse("All services are healthy!", gin.H{
		"status": "healthy",
	})
	c.JSON(http.StatusOK, response)
}

func getServicesHealthStatus(c *gin.Context, h *HealthHandler) bool {
	ctx := c.Request.Context()

	if err := h.authClient.Health(ctx); err != nil {
		return false
	}

	if err := h.videoCatalogClient.Health(ctx); err != nil {
		return false
	}

	if err := h.uploadClient.Health(ctx); err != nil {
		return false
	}

	return true
}
