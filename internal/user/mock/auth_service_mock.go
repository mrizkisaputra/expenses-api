package mock

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (a *AuthServiceMock) Register(ctx context.Context, user *model.User) (*dto.UserResponse, error) {
	args := a.Called(ctx, user)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.UserResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (a *AuthServiceMock) Login(ctx context.Context, user *model.User) (*dto.JwtToken, error) {
	args := a.Called(ctx, user)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.JwtToken), args.Error(1)
	}
	return nil, args.Error(1)
}
