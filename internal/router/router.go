package router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/handler"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/middleware"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/health", handler.Health)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.GET("/profile", middleware.VerifyTokenMiddleware(), handler.Profile)
		auth.POST("/logout", middleware.VerifyTokenMiddleware(), handler.Logout)
	}

	videos := r.Group("/videos")
	{
		videos.GET("", handler.FindAll)
		videos.GET("/:id", handler.FindById)
		videos.POST("/upload/presigned-url", middleware.VerifyTokenMiddleware(), handler.CreatePresignedUrl)
		videos.POST("/upload/webhook", middleware.VerifyTokenMiddleware(), handler.UploadedWebhook)
	}
}
