package converter

import (
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
)

// converter user to UserResponse
func ToUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		Id:          user.Id,
		FirstName:   user.Information.FirstName,
		LastName:    user.Information.LastName,
		Email:       user.Email,
		Password:    user.Password,
		Avatar:      user.Avatar,
		City:        user.Information.City,
		PhoneNumber: user.Information.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
