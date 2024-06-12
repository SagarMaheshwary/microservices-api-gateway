package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	arpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	apb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

func Register(c *gin.Context) {
	in := new(types.RegisterInput)
	ve := new(types.RegisterValidationError)

	if err := c.ShouldBindJSON(&in); err != nil {
		response := helper.PrepareResponseFromValidationError(err, ve)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	response, err := arpc.Auth.Register(&apb.RegisterRequest{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	})

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, ve)
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	in := new(types.LoginInput)
	ve := new(types.LoginValidationError)

	if err := c.ShouldBind(&in); err != nil {
		response := helper.PrepareResponseFromValidationError(err, ve)
		c.JSON(http.StatusBadRequest, response)

		return
	}

	response, err := arpc.Auth.Login(&apb.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	})

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, ve)
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}

func Logout(c *gin.Context) {
	h := new(types.AuthorizationHeader)

	c.ShouldBindHeader(&h)

	response, err := arpc.Auth.Logout(&apb.LogoutRequest{}, h.Token)

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, &types.LogoutValidationError{})
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}
