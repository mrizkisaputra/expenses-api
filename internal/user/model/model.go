package model

import "github.com/google/uuid"

// Model DTO create new user
type UserRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100,alpha"`
	LastName  string `json:"last_name" validate:"required,max=100,alpha"`
	Email     string `json:"email" validate:"required,max=100,email"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}

// UserLoginRequest Model DTO login new user
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,max=100,email"`
	Password string `json:"password" validate:"required,max=100,min=8"`
}

// UserUpdateRequest Model DTO update user
type UserUpdateRequest struct {
	FirstName   string `json:"first_name" validate:"required,max=100,alpha"`
	LastName    string `json:"last_name" validate:"required,max=100,alpha"`
	Email       string `json:"email" validate:"required,max=100,email"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
	Avatar      string `json:"avatar" validate:"url,max=512"`
	City        string `json:"city" validate:"max=100,alpha"`
	PhoneNumber string `json:"phone_number" validate:"max=13,numeric"`
}

// UserResponse Model DTO response user
type UserResponse struct {
	Id          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Avatar      string    `json:"avatar"`
	City        string    `json:"city"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
}

type UserTokenResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Jwt     JwtToken `json:"jwt"`
}

type JwtToken struct {
	AccessToken  string `json:"access_Token"`
	RefreshToken string `json:"refresh_token"`
}

// ApiUserResponse Response API user
type ApiUserResponse struct {
	Status  int
	Message string
	Data    UserResponse
}
