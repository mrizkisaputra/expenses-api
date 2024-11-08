package mock

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/stretchr/testify/mock"
)

// UserPostgresRepositoryMock is a mock type for user.UserPostgresRepository
type UserPostgresRepositoryMock struct {
	mock.Mock
}

func (u *UserPostgresRepositoryMock) Create(ctx context.Context, entity *model.User) (*model.User, error) {
	args := u.Called(ctx, entity)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *UserPostgresRepositoryMock) Update(ctx context.Context, entity *model.User) (*model.User, error) {
	args := u.Called(ctx, entity)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *UserPostgresRepositoryMock) FindByEmail(ctx context.Context, entity *model.User) (*model.User, error) {
	args := u.Called(ctx, entity)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *UserPostgresRepositoryMock) FindById(ctx context.Context, entity *model.User) (*model.User, error) {
	args := u.Called(ctx, entity)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (u *UserPostgresRepositoryMock) FindAlreadyExistByEmail(ctx context.Context, entity *model.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}
