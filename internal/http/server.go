package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/config"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/router"
)

func Connect() {
	c := config.Conf.HTTPServer
	address := fmt.Sprintf("%v:%d", c.Host, c.Port)

	r := gin.Default()

	r.Use(CORSMiddleware()) //@TODO: should it be in code or reverse proxy?

	router.InitRoutes(r)

	s := &http.Server{
		Addr:    address,
		Handler: r,
	}

	log.Info("HTTP server started on %v", address)

	if err := s.ListenAndServe(); err != nil {
		log.Error("HTTP server failed to start %v", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
