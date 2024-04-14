package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cons "github.com/sagarmaheshwary/microservices-api-gateway/internal/constant"
	authrpc "github.com/sagarmaheshwary/microservices-api-gateway/internal/grpc/authentication"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/helper"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/lib/log"
	"github.com/sagarmaheshwary/microservices-api-gateway/internal/types"
	apb "github.com/sagarmaheshwary/microservices-api-gateway/proto/authentication/authentication"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := new(types.AuthorizationHeader)

		if err := ctx.ShouldBindHeader(&h); err != nil {
			response := helper.PrepareResponse(cons.MessageUnauthorized, gin.H{})
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)

			return
		}

		response, err := authrpc.Auth.VerifyToken(&apb.VerifyTokenRequest{}, h.Token)

		if err != nil {
			status, response := helper.PrepareResponseFromgrpcError(err, &types.VerifyTokenValidationError{})
			ctx.AbortWithStatusJSON(status, response)

			return
		}

		log.Info("Token is valid: %v", response)

		ctx.Next()
	}
}
