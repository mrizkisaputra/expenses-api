package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mrizkisaputra/expenses-api/pkg/contextutils"
	"github.com/sirupsen/logrus"
	"time"
)

// RequestLoggerMiddleware is a middleware for logger http request
func (mw *MiddlewareManager) RequestLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		// after request
		latency := time.Since(start)
		mw.logger.WithFields(logrus.Fields{
			"RequestId": contextutils.GetRequestId(ctx),
			"ClientIP":  ctx.ClientIP(),
			"Method":    ctx.Request.Method,
			"UserAgent": ctx.Request.UserAgent(),
			"Path":      ctx.Request.URL.Path,
			"Status":    ctx.Writer.Status(),
			"Latency":   latency,
		}).Info("HTTP Request completed")
	}
}
