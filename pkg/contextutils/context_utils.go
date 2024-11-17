package contextutils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	KeyRequestID = "requestId"
)

// GetRequestId is a function get request id
func GetRequestId(ctx *gin.Context) string {
	requestId := ctx.GetString(KeyRequestID)
	return requestId
}

// GetIPAddress is a function get ip address from request
func GetIPAddress(ctx *gin.Context) string {
	return ctx.ClientIP()
}

func AssignRequestId(ctx *gin.Context) string {
	id := uuid.New().String()

	// save requestId in gin context
	ctx.Set(KeyRequestID, id)

	return id
}
