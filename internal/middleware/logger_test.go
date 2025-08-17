package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
)

func TestZerologMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		route          string
		status         int
		query          string
		expectedFields []string
	}{
		{
			name:   "info log for successful request",
			route:  "/ping",
			status: http.StatusOK,
			query:  "foo=bar",
			expectedFields: []string{
				`"method":"GET"`,
				`"path":"/ping"`,
				`"status":200`,
				`"query":"foo=bar"`,
			},
		},
		{
			name:   "error log for failed request",
			route:  "/fail",
			status: http.StatusBadRequest,
			query:  "",
			expectedFields: []string{
				`"status":400`,
				`"level":"error"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(middleware.ZerologMiddleware())

			switch tt.route {
			case "/ping":
				r.GET("/ping", func(c *gin.Context) {
					c.String(tt.status, "pong")
				})
			case "/fail":
				r.GET("/fail", func(c *gin.Context) {
					c.String(tt.status, "bad request")
				})
			}

			url := tt.route
			if tt.query != "" {
				url += "?" + tt.query
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Assert response status
			assert.Equal(t, tt.status, w.Code, "unexpected response status")

			// Assert log contains expected fields
			logged := buf.String()
			for _, field := range tt.expectedFields {
				assert.Contains(t, logged, field, "log should contain expected field")
			}
		})
	}
}
