package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mrizkisaputra/expenses-api/internal/middleware"
	"github.com/mrizkisaputra/expenses-api/internal/user"
)

// MapAuthRoutes is a function routes
func MapAuthRoutes(authGroup *gin.RouterGroup, controller user.AuthController) {
	authGroup.POST("/register", controller.RegisterNewUser())
	authGroup.POST("/login", controller.LoginNewUser())
}

// MapUserRoutes is a function routes
func MapUserRoutes(userGroup *gin.RouterGroup, controller user.UserController, mw *middleware.MiddlewareManager) {
	userGroup.Use(mw.AuthJwtMiddleware())
	userGroup.GET("/me", controller.GetCurrentUser())
	userGroup.PATCH("/me", controller.Update())
	userGroup.POST("/avatar", controller.PostAvatar())
}
