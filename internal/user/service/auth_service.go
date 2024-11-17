package service

import (
	"context"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/mrizkisaputra/expenses-api/pkg/converter"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/utils"
	"net/http"
)

// authService acts as a struct for injecting an implementation of AuthService interface
// for use in service methods.
type authService struct {
	cfg    *config.Config
	pgRepo user.UserPostgresRepository
}

// NewAuthService is a factory function for
// initializing a authService with its repository layer dependencies
func NewAuthService(config *ServiceConfig) user.AuthService {
	return &authService{
		cfg:    config.Config,
		pgRepo: config.UserPostgresRepository,
	}
}

func (auth *authService) Register(ctx context.Context, user *model.User) (*dto.UserResponse, error) {
	// ensure the email is not already registered
	result, err := auth.pgRepo.FindByEmail(ctx, user)
	if result != nil && err == nil {
		return nil, httpErrors.NewError(http.StatusConflict, httpErrors.EmailAlreadyExistsMsg, err)
	}

	// prepare user for register
	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	// register
	createdUser, err := auth.pgRepo.Create(ctx, user)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}
	return converter.ToUserResponse(createdUser), nil
}

func (auth *authService) Login(ctx context.Context, user *model.User) (*dto.JwtToken, error) {
	// look up requested user
	foundUser, err := auth.pgRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, httpErrors.NewError(http.StatusBadRequest, httpErrors.InvalidEmailOrPasswordMsg, err)
	}

	// compare the request password with the one in the database
	if err := user.ComparePassword(foundUser.Password); err != nil {
		return nil, httpErrors.NewError(http.StatusBadRequest, httpErrors.InvalidEmailOrPasswordMsg, err)
	}

	// generate a JWT token
	accessToken, refreshToken, err := utils.GenerateTokenPair(foundUser, auth.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	return &dto.JwtToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
