package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prachurjya15/simple-bank/token"
)

const (
	authorizationHeaderkey  = "authorization"
	authorizationTypeBearer = "bearer"
	authPayloadKey          = "auth_payload"
)

func AuthMiddleware(tokenMaker token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderkey)
		if len(authHeader) == 0 {
			err := fmt.Errorf("Auth Header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err := fmt.Errorf("Invalid auth format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer {
			err := fmt.Errorf("Not using bearer token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		}
		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
