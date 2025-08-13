package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	authpb "github.com/sagarmaheshwary/microservices-api-gateway/internal/proto/authentication/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
)

type VerifyTokenFunc func(ctx context.Context, in *authpb.VerifyTokenRequest, token string) (*authpb.VerifyTokenResponse, error)

func VerifyTokenMiddleware(verifyToken VerifyTokenFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var h types.AuthorizationHeader
		if err := c.ShouldBindHeader(&h); err != nil {
			response := helper.PrepareResponse(constant.MessageUnauthorized, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		response, err := verifyToken(c.Request.Context(), &authpb.VerifyTokenRequest{}, h.Token)
		if err != nil {
			status, response := helper.PrepareResponseFromGRPCError(err, &types.VerifyTokenValidationError{})
			c.AbortWithStatusJSON(status, response)
			return
		}

		c.Set(constant.AuthUser, response.Data.User)
		c.Set(constant.GRPCHeaderAuthorization, h)
		c.Next()
	}
}
