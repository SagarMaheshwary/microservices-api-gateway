package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	cons "github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	urpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	apb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	upb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

func CreatePresignedUrl(c *gin.Context) {
	response, err := urpc.Upload.CreatePresignedUrl(&upb.CreatePresignedUrlRequest{})

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, gin.H{})
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}

func UploadedWebhook(c *gin.Context) {
	in := new(types.UploadedWebhookInput)
	ve := new(types.UploadedWebhookValidationError)

	u, exists := c.Get(cons.AuthUser)

	if !exists {
		log.Error("Authenticated user does not exists in context!")

		c.JSON(http.StatusInternalServerError, helper.PrepareResponse(cons.MessageInternalServerError, gin.H{}))
		return
	}

	user := u.(*apb.User)

	if err := c.ShouldBind(&in); err != nil {
		response := helper.PrepareResponseFromValidationError(err, ve)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	response, err := urpc.Upload.UploadedWebhook(&upb.UploadedWebhookRequest{
		VideoId:     in.VideoId,
		ThumbnailId: in.ThumbnailId,
		Title:       in.Title,
		Description: in.Description,
	}, strconv.Itoa(int(user.Id)))

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, ve)
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}
