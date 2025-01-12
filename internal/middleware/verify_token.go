package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := new(types.AuthorizationHeader)

		if err := c.ShouldBindHeader(&h); err != nil {
			response := helper.PrepareResponse(constant.MessageUnauthorized, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)

			return
		}

		response, err := authrpc.Auth.VerifyToken(&authpb.VerifyTokenRequest{}, h.Token)

		if err != nil {
			status, response := helper.PrepareResponseFromgrpcError(err, &types.VerifyTokenValidationError{})
			c.AbortWithStatusJSON(status, response)

			return
		}

		c.Set(constant.AuthUser, response.Data.User)

		c.Next()
	}
}
