package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	uploadrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/upload"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	uploadpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/upload/upload"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

type UploadHandler struct {
	uploadClient uploadrpc.UploadService
}

func NewUploadHandler(c uploadrpc.UploadService) *UploadHandler {
	return &UploadHandler{uploadClient: c}
}

func (u *UploadHandler) CreatePresignedUrl(c *gin.Context) {
	res, err := u.uploadClient.CreatePresignedUrl(c.Request.Context(), &uploadpb.CreatePresignedUrlRequest{})
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, gin.H{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (u *UploadHandler) UploadedWebhook(c *gin.Context) {
	user, exists := c.Get(constant.AuthUser)
	if !exists {
		logger.Error("Authenticated user does not exists in context!")
		c.JSON(http.StatusInternalServerError, helper.PrepareResponse(constant.MessageInternalServerError, gin.H{}))
		return
	}
	authUser := user.(*authpb.User)

	var in types.UploadedWebhookInput
	if err := c.ShouldBind(&in); err != nil {
		res := helper.PrepareResponseFromValidationError(err, &types.UploadedWebhookValidationError{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	req := &uploadpb.UploadedWebhookRequest{
		VideoId:     in.VideoId,
		ThumbnailId: in.ThumbnailId,
		Title:       in.Title,
		Description: in.Description,
	}

	res, err := u.uploadClient.UploadedWebhook(c.Request.Context(), req, strconv.Itoa(int(authUser.Id)))
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, &types.UploadedWebhookValidationError{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
