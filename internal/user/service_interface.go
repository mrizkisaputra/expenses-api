package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
)

// AuthService defines methods the layer controller expects.
// any services it interacts with to implement.
type AuthService interface {
	Register(ctx context.Context, user *model.User) (*dto.UserResponse, error)
	Login(ctx context.Context, user *model.User) (*dto.JwtToken, error)
}

// UserService defines methods the layer controller expects.
// any services it interacts with to implement.
type UserService interface {
	GetCurrentUser(ctx context.Context, id string) (*dto.UserResponse, error)

	Update(ctx context.Context, user *model.User) (*dto.UserResponse, error)

	UploadAvatar(ctx context.Context, file *model.UserUploadInput, id uuid.UUID) (*dto.UserResponse, error)
}
