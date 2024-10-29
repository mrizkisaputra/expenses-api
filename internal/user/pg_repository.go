package user

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
)

// User Postgres repository interface
type UserPostgresRepository interface {
	Register(ctx context.Context, entity *model.User) (*model.User, error)

	Update(ctx context.Context, entity *model.User) (*model.User, error)

	FindByEmail(ctx context.Context, entity *model.User) (*model.User, error)

	FindById(ctx context.Context, entity *model.User) (*model.User, error)

	FindAlreadyExistByEmail(ctx context.Context, entity *model.User) (int64, error)
}
