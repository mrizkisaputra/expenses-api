package mock

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/stretchr/testify/mock"
	"time"
)

type UserRedisRepositoryMock struct {
	mock.Mock
}

func (u *UserRedisRepositoryMock) Set(ctx context.Context, key string, expiration time.Duration, value *model.User) error {
	args := u.Called(ctx, key, expiration, value)
	return args.Error(0)
}

func (u *UserRedisRepositoryMock) Get(ctx context.Context, key string) (*model.User, error) {
	args := u.Called(ctx, key)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *UserRedisRepositoryMock) Delete(ctx context.Context, key string) error {
	args := u.Called(ctx, key)
	return args.Error(0)
}
