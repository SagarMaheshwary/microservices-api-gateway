package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/prometheus"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
)

func resetMetrics() {
	prometheus.TotalRequests.Reset()
	prometheus.ErrorCount.Reset()
	prometheus.RequestDuration.Reset()
	prometheus.ServiceHealth.Set(0)
}

func TestPrometheusMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		status         int
		expectedRoute  string
		expectErrorInc bool
		registerRoute  bool
	}{
		{
			name:           "successful request increments totalRequests and observes duration",
			method:         http.MethodGet,
			path:           "/test",
			status:         http.StatusOK,
			expectedRoute:  "/test",
			expectErrorInc: false,
			registerRoute:  true,
		},
		{
			name:           "failed request increments error counter",
			method:         http.MethodPost,
			path:           "/error",
			status:         http.StatusInternalServerError,
			expectedRoute:  "/error",
			expectErrorInc: true,
			registerRoute:  true,
		},
		{
			name:           "request with no defined route sets unknown label",
			method:         http.MethodGet,
			path:           "/notfound",
			status:         http.StatusNotFound,
			expectedRoute:  "unknown",
			expectErrorInc: true,
			registerRoute:  false, // don't register, so FullPath() == "" -> "unknown"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetMetrics()

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(middleware.PrometheusMiddleware())

			if tt.registerRoute {
				switch tt.path {
				case "/test":
					r.GET("/test", func(c *gin.Context) { c.Status(tt.status) })
				case "/error":
					r.POST("/error", func(c *gin.Context) { c.Status(tt.status) })
				}
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			statusLabel := strconv.Itoa(tt.status)

			// Verify TotalRequests incremented
			gotTotal := testutil.ToFloat64(prometheus.TotalRequests.WithLabelValues(tt.method, tt.expectedRoute, statusLabel))
			assert.Equal(t, float64(1), gotTotal, "TotalRequests should be incremented once")

			// Verify RequestDuration observed
			gotDurationObs := testutil.CollectAndCount(prometheus.RequestDuration)
			require.NotZero(t, gotDurationObs, "RequestDuration should have at least one observation")

			// Verify ErrorCount increments based on status
			gotError := testutil.ToFloat64(prometheus.ErrorCount.WithLabelValues(tt.method, tt.expectedRoute))
			if tt.expectErrorInc {
				assert.Equal(t, float64(1), gotError, "ErrorCount should be incremented")
			} else {
				assert.Equal(t, float64(0), gotError, "ErrorCount should not be incremented")
			}
		})
	}
}
