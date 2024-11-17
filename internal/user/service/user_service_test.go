package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/mrizkisaputra/expenses-api/config"
	mockObject "github.com/mrizkisaputra/expenses-api/internal/user/mock"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/pkg/converter"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/logger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

var cfg = &config.Config{
	AWS: config.AwsConfig{
		MinioEndpoint: "http://localhost:9000",
	},
}

func TestUserService_GetCurrentUser(t *testing.T) {
	// scenario test case #1
	t.Run("[Test Case #1] Should return an error when an invalid ID is provided", func(t *testing.T) {
		// create instance service
		us := NewUserService(&ServiceConfig{})

		invalidId := "invalid-id"
		_, err := us.GetCurrentUser(context.Background(), invalidId)

		require.Error(t, err)
		require.NotNil(t, err)
		var er *httpErrors.Error
		if errors.As(err, &er) {
			require.Equal(t, http.StatusInternalServerError, er.Status)
		}
	})

	// scenario test case #2
	t.Run("[Test Case #2] Successfully get user from Redis cache", func(t *testing.T) {
		// create instance service and mock object
		mockRedisRepo := new(mockObject.UserRedisRepositoryMock)
		us := NewUserService(&ServiceConfig{
			UserRedisRepository: mockRedisRepo,
		})

		mockUser := &model.User{
			Email:    "testcase@test.com",
			Password: "secret",
			Information: model.Information{
				FirstName:   "test",
				LastName:    "test",
				City:        nil,
				PhoneNumber: nil,
			},
			Avatar:    nil,
			UpdatedAt: time.Now().UnixMilli(),
			CreatedAt: time.Now().UnixMilli(),
		}
		mockRedisRepo.On("Get", mock.Anything, mock.Anything).Return(mockUser, nil)

		validId := uuid.New()
		response, err := us.GetCurrentUser(context.Background(), validId.String())

		require.Nil(t, err)
		require.Equal(t, converter.ToUserResponse(mockUser), response)
		mockRedisRepo.AssertExpectations(t)
	})

	// scenario test case #3
	t.Run("[Test Case #3] Successfully get user from Postgres database and cache in Redis", func(t *testing.T) {
		// create instance service and mock object
		mockPgRepo := new(mockObject.UserPostgresRepositoryMock)
		mockRedisRepo := new(mockObject.UserRedisRepositoryMock)
		us := NewUserService(&ServiceConfig{
			UserPostgresRepository: mockPgRepo,
			UserRedisRepository:    mockRedisRepo,
			Logger:                 logger.NewLogrusLogger(cfg),
		})

		validId := uuid.New()
		mockUser := &model.User{
			Id:       validId,
			Email:    "testcase@test.com",
			Password: "secret",
			Information: model.Information{
				FirstName:   "test",
				LastName:    "test",
				City:        nil,
				PhoneNumber: nil,
			},
			Avatar:    nil,
			UpdatedAt: time.Now().UnixMilli(),
			CreatedAt: time.Now().UnixMilli(),
		}
		mockRedisRepo.On("Get", mock.Anything, mock.Anything).Return(nil, nil)
		mockPgRepo.On("FindById", mock.Anything, mock.Anything).Return(mockUser, nil)
		mockRedisRepo.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		response, err := us.GetCurrentUser(context.Background(), validId.String())

		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, converter.ToUserResponse(mockUser), response)
		mockPgRepo.AssertExpectations(t)
		mockRedisRepo.AssertExpectations(t)
	})

	// scenario test case #4
	t.Run("[Test Case #4] Error when user not found in Postgres database", func(t *testing.T) {
		// create instance service and mock object
		mockPgRepo := new(mockObject.UserPostgresRepositoryMock)
		mockRedisRepo := new(mockObject.UserRedisRepositoryMock)
		us := NewUserService(&ServiceConfig{
			UserPostgresRepository: mockPgRepo,
			UserRedisRepository:    mockRedisRepo,
			Logger:                 logger.NewLogrusLogger(cfg),
		})

		mockRedisRepo.On("Get", mock.Anything, mock.Anything).Return(nil, nil)
		mockPgRepo.On("FindById", mock.Anything, mock.Anything).Return(nil, httpErrors.NewNotFoundError(nil))
		mockRedisRepo.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		validId := uuid.New()
		response, err := us.GetCurrentUser(context.Background(), validId.String())

		require.Error(t, err)
		require.NotNil(t, err)
		require.Nil(t, response)

		var er httpErrors.Error
		if errors.As(err, &er) {
			require.Equal(t, http.StatusNotFound, er.Status)
		}
	})
}

func TestUserService_Update(t *testing.T) {
	// scenario test case #1
	t.Run("[Test Case #1] Successfully updated data user and delete cache in redis", func(t *testing.T) {
		mockPgRepo := new(mockObject.UserPostgresRepositoryMock)
		mockRedisRepo := new(mockObject.UserRedisRepositoryMock)
		us := NewUserService(&ServiceConfig{
			UserPostgresRepository: mockPgRepo,
			UserRedisRepository:    mockRedisRepo,
		})

		userID := uuid.New()
		city := "south sumatera"

		// Mock data pengguna lama
		oldUser := &model.User{
			Id:    userID,
			Email: "oldemail@example.com",
			Information: model.Information{
				FirstName:   "OldFirstName",
				LastName:    "OldLastName",
				City:        &city,
				PhoneNumber: nil,
			},
		}

		// Mock data pengguna baru untuk diupdate
		updateUser := &model.User{
			Id:    userID,
			Email: "newemail@example.com",
			Information: model.Information{
				FirstName:   "NewFirstName",
				LastName:    "",  // Tidak diupdate karena kosong
				City:        nil, // Tidak diupdate karena nil
				PhoneNumber: nil, // Tidak diupdate karena nil
			},
		}

		mockPgRepo.On("FindById", mock.Anything, mock.Anything).Return(oldUser, nil)
		mockPgRepo.On("Update", mock.Anything, mock.Anything).Return(oldUser, nil)
		mockRedisRepo.On("Delete", mock.Anything, mock.Anything).Return(nil)

		response, err := us.Update(context.Background(), updateUser)

		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, "newemail@example.com", response.Email)
		require.Equal(t, "NewFirstName", response.FirstName)
		mockPgRepo.AssertExpectations(t)
		mockRedisRepo.AssertExpectations(t)
	})

	// scenario test case #2
	t.Run("[Test Case #2] Error when update user", func(t *testing.T) {
		mockPgRepo := new(mockObject.UserPostgresRepositoryMock)
		us := NewUserService(&ServiceConfig{UserPostgresRepository: mockPgRepo})

		id := uuid.New()
		mockUser := &model.User{
			Id:    id,
			Email: "mrizkisaputra_updated@gmail.com",
			Information: model.Information{
				FirstName: "Muhammat updated",
				LastName:  "Saputra updated",
			},
		}

		mockPgRepo.On("FindById", mock.Anything, mock.Anything).Return(mockUser, nil)
		mockPgRepo.On("Update", mock.Anything, mock.Anything).Return(nil, httpErrors.NewInternalServerError(nil))

		response, err := us.Update(context.Background(), mockUser)

		require.Nil(t, response)
		require.Error(t, err)
		require.NotNil(t, err)

		var er *httpErrors.Error
		if errors.As(err, &er) {
			require.Equal(t, http.StatusInternalServerError, er.Status)
		}

		mockPgRepo.AssertExpectations(t)
	})

	// scenario test case #3

}

func TestUserService_UploadAvatar(t *testing.T) {
	t.Run("Successfully upload avatar", func(t *testing.T) {
		mockPgRepo := new(mockObject.UserPostgresRepositoryMock)
		mockAwsRepo := new(mockObject.AwsRepositoryMock)
		us := NewUserService(&ServiceConfig{
			UserPostgresRepository: mockPgRepo,
			AwsUserRepository:      mockAwsRepo,
			Config:                 cfg,
		})

		mockFile := model.UserUploadInput{}
		uploadInfo := &minio.UploadInfo{}

		id := uuid.New()
		mockUser := &model.User{
			Id:       id,
			Email:    "mrizkisaputra@gmail.com",
			Password: "1234",
		}

		// arrange mock
		mockAwsRepo.On("PutObject", mock.Anything, mock.Anything).Return(uploadInfo, nil)
		mockPgRepo.On("Update", mock.Anything, mock.Anything).Return(mockUser, nil)

		response, err := us.UploadAvatar(context.Background(), id, &mockFile)
		require.Nil(t, err)
		require.NoError(t, err)
		require.NotNil(t, response)

		mockPgRepo.AssertExpectations(t)
		mockAwsRepo.AssertExpectations(t)
	})
}
