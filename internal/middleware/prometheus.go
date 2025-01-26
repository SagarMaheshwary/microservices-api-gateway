package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		route := c.FullPath()
		if route == "" {
			route = "unknown" // Handle requests with no defined route.
		}

		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())

		prometheus.TotalRequests.WithLabelValues(method, route, status).Inc()
		prometheus.RequestDuration.WithLabelValues(method, route).Observe(duration)

		if c.Writer.Status() >= 400 {
			prometheus.ErrorCount.WithLabelValues(method, route).Inc()
		}
	}
}
