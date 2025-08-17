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
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.PrepareResponse(constant.MessageUnauthorized, gin.H{}))
			return
		}

		res, err := verifyToken(c.Request.Context(), &authpb.VerifyTokenRequest{}, h.Token)
		if err != nil {
			status, res := helper.PrepareResponseFromGRPCError(err, &types.VerifyTokenValidationError{})
			c.AbortWithStatusJSON(status, res)
			return
		}

		c.Set(constant.AuthUser, res.Data.User)
		c.Set(constant.GRPCHeaderAuthorization, h)
		c.Next()
	}
}
