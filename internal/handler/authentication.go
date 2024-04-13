package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cons "github.com/sagarmaheshwary/microservices-api-gateway/internal/constants"
	auth "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
	apb "github.com/sagarmaheshwary/microservices-api-gateway/proto/authentication/authentication"
)

func Register(c *gin.Context) {
	in := new(types.RegisterInput)
	ve := new(types.RegisterValidationError)

	if err := c.ShouldBindJSON(&in); err != nil {
		errors := helper.TransformValidationErrors(err, ve)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": cons.MSGBadRequest,
			"data": gin.H{
				"errors": errors,
			},
		})

		return
	}

	response, err := auth.Auth.Register(&apb.RegisterRequest{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	})

	if err != nil {
		status, res := helper.PrepareResponseFromgrpcError(err, ve)
		c.JSON(status, res)

		return
	}

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	in := new(types.LoginInput)
	ve := new(types.LoginValidationError)

	if err := c.ShouldBind(&in); err != nil {
		errors := helper.TransformValidationErrors(err, ve)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": cons.MSGBadRequest,
			"data": gin.H{
				"errors": errors,
			},
		})

		return
	}

	response, err := auth.Auth.Login(&apb.LoginRequest{
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

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": cons.MSGUnauthorized,
			"data":    gin.H{},
		})

		return
	}

	response, err := auth.Auth.Logout(&apb.LogoutRequest{}, h.Token)

	if err != nil {
		status, response := helper.PrepareResponseFromgrpcError(err, &types.LogoutValidationError{})
		c.JSON(status, response)

		return
	}

	c.JSON(http.StatusOK, response)
}
