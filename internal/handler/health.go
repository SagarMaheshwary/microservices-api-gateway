package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
)

func Health(c *gin.Context) {
	if !getServicesHealthStatus() {
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

func getServicesHealthStatus() bool {
	if !authrpc.HealthCheck() {
		return false
	}

	if !videocatalogrpc.HealthCheck() {
		return false
	}

	if !uploadrpc.HealthCheck() {
		return false
	}

	return true
}
