package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
)

func InitRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/logout", middleware.VerifyToken(), handler.Logout)
	}

	videos := r.Group("/videos")
	{
		videos.POST("/upload/presigned-url", middleware.VerifyToken(), handler.CreatePresignedUrl)
		videos.POST("/upload/webhook", middleware.VerifyToken(), handler.UploadedWebhook)
	}
}
