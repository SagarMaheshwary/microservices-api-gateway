package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
)

func Health(c *gin.Context) {

	if !getServicesHealthStatus() {
		c.JSON(http.StatusServiceUnavailable, helper.PrepareResponse("Some services are not available!", gin.H{
			"status": "degraded",
		}))

		return
	}

	c.JSON(http.StatusOK, helper.PrepareResponse("All services are healthy!", gin.H{
		"status": "healthy",
	}))
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
