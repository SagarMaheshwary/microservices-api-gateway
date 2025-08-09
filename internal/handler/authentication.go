package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

type AuthenticationHandler struct {
	authClient authrpc.AuthenticationService
}

func NewAuthHandler(c authrpc.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{authClient: c}
}

func (a *AuthenticationHandler) Register(c *gin.Context) {
	var in types.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, helper.PrepareResponseFromValidationError(err, &types.RegisterValidationError{}))
		return
	}

	req := &authpb.RegisterRequest{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}

	res, err := a.authClient.Register(c.Request.Context(), req)
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, &types.RegisterValidationError{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (a *AuthenticationHandler) Login(c *gin.Context) {
	var in types.LoginInput
	if err := c.ShouldBind(&in); err != nil {
		res := helper.PrepareResponseFromValidationError(err, &types.LoginValidationError{})
		c.JSON(http.StatusBadRequest, res)
		return
	}

	req := &authpb.LoginRequest{
		Email:    in.Email,
		Password: in.Password,
	}

	res, err := a.authClient.Login(c.Request.Context(), req)
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, &types.LoginValidationError{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *AuthenticationHandler) Profile(c *gin.Context) {
	user, _ := c.Get(constant.AuthUser)

	res := helper.PrepareResponse(constant.MessageOK, gin.H{
		constant.AuthUser: user,
	})
	c.JSON(http.StatusOK, res)
}

func (a *AuthenticationHandler) Logout(c *gin.Context) {
	h := new(types.AuthorizationHeader)
	c.ShouldBindHeader(&h)

	res, err := a.authClient.Logout(c.Request.Context(), &authpb.LogoutRequest{}, h.Token)
	if err != nil {
		status, res := helper.PrepareResponseFromGRPCError(err, &types.LogoutValidationError{})
		c.JSON(status, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
