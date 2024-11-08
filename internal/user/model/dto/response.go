package dto

import "github.com/google/uuid"

type ApiUserResponse struct {
	Status  int
	Message string
	Data    interface{}
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

// UserResponse Model DTO response user
type UserResponse struct {
	Id          uuid.UUID `json:"id,omitempty"`
	FirstName   string    `json:"first_name,omitempty"`
	LastName    string    `json:"last_name,omitempty"`
	Email       string    `json:"email,omitempty"`
	Password    string    `json:"password,omitempty"`
	Avatar      *string   `json:"avatar,omitempty"`
	City        *string   `json:"city,omitempty"`
	PhoneNumber *string   `json:"phone_number,omitempty"`
	CreatedAt   int64     `json:"created_at,omitempty"`
	UpdatedAt   int64     `json:"updated_at,omitempty"`
}
