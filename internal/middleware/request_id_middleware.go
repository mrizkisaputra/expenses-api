package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mrizkisaputra/expenses-api/pkg/contextutils"
)

// RequestIdMiddleware is a middleware for assign request id
func (mw *MiddlewareManager) RequestIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// before request
		requestId := ctx.GetHeader("X-Request-Id")
		if requestId == "" {
			requestId = contextutils.AssignRequestId(ctx)
		}

		// added requestId in header response
		ctx.Writer.Header().Set("X-Request-Id", requestId)

		ctx.Next()
	}
}
