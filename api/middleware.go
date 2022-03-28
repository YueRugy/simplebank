package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simplebank/token"
	"strings"
)

const (
	authorizationHeadKey    = "authorization"
	authorizationBearer     = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

var (
	authAccountNotBelongUserError =errors.New("account doesn't belong authorized user")
)

func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHead := ctx.GetHeader(authorizationHeadKey)
		if len(authorizationHead) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}

		fields := strings.Fields(authorizationHead)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationBearer {
			err := fmt.Errorf("unsupport authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, responseError(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
