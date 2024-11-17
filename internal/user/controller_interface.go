package user

import "github.com/gin-gonic/gin"

// AuthController defines methods the routes expects
// any controllers it interacts with to implement
type AuthController interface {
	RegisterNewUser() gin.HandlerFunc
	LoginNewUser() gin.HandlerFunc
}

// UserController defines methods the routes expects
// any controllers it interacts with to implement
type UserController interface {
	GetCurrentUser() gin.HandlerFunc
	Update() gin.HandlerFunc
	PostAvatar() gin.HandlerFunc
}
