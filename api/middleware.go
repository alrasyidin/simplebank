package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alrasyidin/simplebank-go/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey     = "authorization"
	authorizationHeaderType    = "bearer"
	authorizationHeaderPayload = "authorization_payload"
)

func authMiddleware(generator token.Generator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization := ctx.GetHeader(authorizationHeaderKey)

		if len(authorization) == 0 {
			err := errors.New("authorization headers isn't provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorization)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != authorizationHeaderType {
			err := errors.New("unsupported authorization header key")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		accessToken := fields[1]
		payload, err := generator.VerifyToken(accessToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationHeaderPayload, payload)
		ctx.Next()
	}
}
