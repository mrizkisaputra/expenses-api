package user

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"time"
)

type UserRedisRepository interface {
	Set(ctx context.Context, key string, expiration time.Duration, value *model.User) error

	Get(ctx context.Context, key string) (*model.User, error)

	Delete(ctx context.Context, key string) error
}
