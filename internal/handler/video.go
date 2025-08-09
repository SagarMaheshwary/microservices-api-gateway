package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	videocatalogrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/video-catalog"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/logger"
	videocatalogpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/video_catalog"
)

type VideoCatalogHandler struct {
	videoCatalogClient videocatalogrpc.VideoCatalogService
}

func NewVideoCatalogHandler(c videocatalogrpc.VideoCatalogService) *VideoCatalogHandler {
	return &VideoCatalogHandler{videoCatalogClient: c}
}

func (v *VideoCatalogHandler) FindAll(c *gin.Context) {
	res, err := v.videoCatalogClient.FindAll(c.Request.Context(), &videocatalogpb.FindAllRequest{})
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, gin.H{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (v *VideoCatalogHandler) FindById(c *gin.Context) {
	videoId := c.Param("id")
	id, err := strconv.Atoi(videoId)
	if err != nil {
		logger.Error("Unable to parse video id %v", err)
		c.JSON(http.StatusBadRequest, helper.PrepareResponse(constant.MessageBadRequest, gin.H{}))
		return
	}

	req := &videocatalogpb.FindByIdRequest{
		Id: int32(id),
	}

	res, err := v.videoCatalogClient.FindById(c.Request.Context(), req)
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, gin.H{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
