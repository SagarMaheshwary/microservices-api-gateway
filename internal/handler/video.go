package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

func CreatePresignedUrl(c *gin.Context) {
	response, err := uploadrpc.Upload.CreatePresignedUrl(c.Request.Context(), &uploadpb.CreatePresignedUrlRequest{})
	if err != nil {
		status, response := helper.PrepareResponseFromGrpcError(err, gin.H{})
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func UploadedWebhook(c *gin.Context) {
	u, exists := c.Get(constant.AuthUser)
	if !exists {
		logger.Error("Authenticated user does not exists in context!")
		c.JSON(http.StatusInternalServerError, helper.PrepareResponse(constant.MessageInternalServerError, gin.H{}))
		return
	}
	user := u.(*authpb.User)

	in := new(types.UploadedWebhookInput)
	ve := new(types.UploadedWebhookValidationError)
	if err := c.ShouldBind(&in); err != nil {
		response := helper.PrepareResponseFromValidationError(err, ve)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response, err := uploadrpc.Upload.UploadedWebhook(
		c.Request.Context(),
		&uploadpb.UploadedWebhookRequest{
			VideoId:     in.VideoId,
			ThumbnailId: in.ThumbnailId,
			Title:       in.Title,
			Description: in.Description,
		},
		strconv.Itoa(int(user.Id)),
	)
	if err != nil {
		status, response := helper.PrepareResponseFromGrpcError(err, ve)
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func FindAll(c *gin.Context) {
	response, err := videocatalogrpc.VideoCatalog.FindAll(c.Request.Context(), &videocatalogpb.FindAllRequest{})
	if err != nil {
		status, response := helper.PrepareResponseFromGrpcError(err, gin.H{})
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

func FindById(c *gin.Context) {
	videoId := c.Param("id")
	id, err := strconv.Atoi(videoId)
	if err != nil {
		logger.Error("Unable to parse video id %v", err)
		c.JSON(http.StatusInternalServerError, helper.PrepareResponse(constant.MessageInternalServerError, gin.H{}))
		return
	}

	response, err := videocatalogrpc.VideoCatalog.FindById(c.Request.Context(), &videocatalogpb.FindByIdRequest{
		Id: int32(id),
	})
	if err != nil {
		status, response := helper.PrepareResponseFromGrpcError(err, gin.H{})
		c.JSON(status, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
