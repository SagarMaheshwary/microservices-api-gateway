package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
)

func InitRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/logout", handler.Logout)
	}

	r.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"test": "test"}) })
}
