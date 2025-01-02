package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	urpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	vcrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	apb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	upb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	vcpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
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

	u, exists := c.Get(constant.AuthUser)

	if !exists {
		log.Error("Authenticated user does not exists in context!")

		c.JSON(http.StatusInternalServerError, helper.PrepareResponse(constant.MessageInternalServerError, gin.H{}))
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

func FindAll(c *gin.Context) {
	response, err := vcrpc.VideoCatalog.FindAll(&vcpb.FindAllRequest{})

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, gin.H{})
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}

func FindById(c *gin.Context) {
	videoId := c.Param("id")

	id, err := strconv.Atoi(videoId)

	if err != nil {
		log.Error("Unable to parse video id %v", err)
		c.JSON(http.StatusInternalServerError, helper.PrepareResponse(constant.MessageInternalServerError, gin.H{}))

		return
	}

	response, err := vcrpc.VideoCatalog.FindById(&vcpb.FindByIdRequest{
		Id: int32(id),
	})

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, gin.H{})
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}
