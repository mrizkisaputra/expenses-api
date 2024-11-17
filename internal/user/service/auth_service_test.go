package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/config"
	mockObject "github.com/mrizkisaputra/expenses-api/internal/user/mock"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestAuthService_Register_Failed(t *testing.T) {
	// instance mock object & authService
	var mockRepo = new(mockObject.UserPostgresRepositoryMock)
	var auth = NewAuthService(&ServiceConfig{UserPostgresRepository: mockRepo})

	mockExistingUser := &model.User{
		Email: "emailalreadyexist@gmail.com",
	}

	// arrange mock for failure case
	mockRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(mockExistingUser, nil)
	response, err := auth.Register(context.Background(), mockExistingUser)
	require.Error(t, err)
	require.Nil(t, response)

	var er httpErrors.Error
	if errors.As(err, &er) {
		require.Equal(t, http.StatusConflict, er.Status)
		require.Equal(t, httpErrors.EmailAlreadyExistsMsg, er.Message)
	}

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	// instance mock object & authService
	var mockRepo = new(mockObject.UserPostgresRepositoryMock)
	var auth = NewAuthService(&ServiceConfig{UserPostgresRepository: mockRepo})

	expectedResponse := &dto.UserResponse{
		Id:        uuid.New(),
		Email:     "mrizkisaputra@test.com",
		Password:  "secret",
		FirstName: "muhammat",
		LastName:  "saputra",
		UpdatedAt: time.Now().UnixMilli(),
		CreatedAt: time.Now().UnixMilli(),
	}

	// initializing mock data
	mockNewUser := &model.User{
		Email:    "mrizkisaputra@test.com",
		Password: "secret",
		Information: model.Information{
			FirstName: "muhammat",
			LastName:  "saputra",
		},
	}

	// arrange mock for successfully case
	mockRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New(httpErrors.EmailAlreadyExistsMsg))
	mockRepo.On("Create", mock.Anything, mock.Anything).
		Return(&model.User{
			Id:       expectedResponse.Id,
			Email:    "mrizkisaputra@test.com",
			Password: "secret",
			Information: model.Information{
				FirstName: "muhammat",
				LastName:  "saputra",
			},
			UpdatedAt: expectedResponse.UpdatedAt,
			CreatedAt: expectedResponse.CreatedAt,
		}, nil)

	response, err := auth.Register(context.Background(), mockNewUser)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, expectedResponse, response)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(mockObject.UserPostgresRepositoryMock)
	auth := NewAuthService(&ServiceConfig{
		Config: &config.Config{
			Server: config.ServerConfig{
				JWTSecretKey: "secret_key",
			},
		},
		UserPostgresRepository: mockRepo,
	})

	// arrange mock for successfully login
	mockRepo.On("FindByEmail", mock.Anything, mock.Anything).
		Return(&model.User{
			Email:    "mrizkisaputra@test.com",
			Password: "$2a$12$WX/or6PO4Ue2CqmlOBMREufvVDcDuSq9cx/AgcqYU/pf.lkiTEP1O", //secret
		}, nil)

	response, err := auth.Login(context.Background(), &model.User{
		Email:    "mrizkisaputra@test.com",
		Password: "secret",
	})

	require.Nil(t, err)
	require.NoError(t, err)
	require.NotNil(t, response)
	//fmt.Printf("accessToken: %s\nrefreshToken: %s\n", response.AccessToken, response.RefreshToken)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_EmailInvalid(t *testing.T) {
	mockRepo := new(mockObject.UserPostgresRepositoryMock)
	auth := NewAuthService(&ServiceConfig{
		Config: &config.Config{
			Server: config.ServerConfig{
				JWTSecretKey: "secret_key",
			},
		},
		UserPostgresRepository: mockRepo,
	})

	// arrange mock for failed login
	mockRepo.On("FindByEmail", mock.Anything, mock.Anything).Return(
		nil,
		httpErrors.NewError(http.StatusBadRequest, httpErrors.InvalidEmailOrPasswordMsg, nil),
	)

	response, err := auth.Login(context.Background(), &model.User{
		Email:    "mrizkisaputra@test.com",
		Password: "secret",
	})

	require.Error(t, err)
	require.NotNil(t, err)

	var er httpErrors.Error
	if errors.As(err, &er) {
		require.Equal(t, http.StatusBadRequest, er.Status)
		require.Equal(t, httpErrors.InvalidEmailOrPasswordMsg, er.Message)
	}

	require.Nil(t, response)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_PasswordInvalid(t *testing.T) {
	mockRepo := new(mockObject.UserPostgresRepositoryMock)
	auth := NewAuthService(&ServiceConfig{
		Config: &config.Config{
			Server: config.ServerConfig{
				JWTSecretKey: "secret_key",
			},
		},
		UserPostgresRepository: mockRepo,
	})

	// arrange mock for failed login, cause password not match
	mockRepo.On("FindByEmail", mock.Anything, mock.Anything).
		Return(&model.User{
			Email:    "mrizkisaputra@test.com",
			Password: "$2a$12$WX/or6PO4Ue2CqmlOBMREufvVDcDuSq9cx/AgcqYU/pf.lkiTEP1O", //secret
		}, nil)

	response, err := auth.Login(context.Background(), &model.User{
		Email:    "mrizkisaputra@test.com",
		Password: "not matching",
	})

	require.Error(t, err)
	require.NotNil(t, err)
	require.Nil(t, response)

	var er httpErrors.Error
	if errors.As(err, &er) {
		require.Equal(t, http.StatusBadRequest, er.Status)
		require.Equal(t, httpErrors.InvalidEmailOrPasswordMsg, er.Message)
	}

	mockRepo.AssertExpectations(t)
}
