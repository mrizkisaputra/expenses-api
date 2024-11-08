package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/mrizkisaputra/expenses-api/internal/user/model/dto"
	"github.com/mrizkisaputra/expenses-api/pkg/converter"
	"github.com/mrizkisaputra/expenses-api/pkg/httpErrors"
	"github.com/mrizkisaputra/expenses-api/pkg/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	basePrefix    = "user-api"
	cacheDuration = 3600
)

// userService acts as a struct for injecting an implementation of UserService interface
// for use in service methods.
type userService struct {
	cfg       *config.Config
	pgRepo    user.UserPostgresRepository
	redisRepo user.UserRedisRepository
	awsRepo   user.AWSUserRepository
	logger    *logrus.Logger
}

// NewUserService is a factory function for
// initializing a userService with its repository layer dependencies.
func NewUserService(sc *ServiceConfig) user.UserService {
	return &userService{
		cfg:       sc.Config,
		pgRepo:    sc.UserPostgresRepository,
		redisRepo: sc.UserRedisRepository,
		awsRepo:   sc.AwsUserRepository,
		logger:    sc.Logger,
	}
}

func (u *userService) GetCurrentUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "UserService.GetCurrentUser.uuid.Parse"))
	}

	// get data from redis cache
	cacheUser, err := u.redisRepo.Get(ctx, utils.GetRedisKey(basePrefix, id))
	if err != nil {
		u.logger.WithError(err).Error("UserService.GetCurrentUser.redisRepo.Get")
	}

	if cacheUser != nil {
		return converter.ToUserResponse(cacheUser), nil
	}

	// get current user by id from postgres db
	currentUser, err := u.pgRepo.FindById(ctx, &model.User{Id: parsedId})
	if err != nil {
		return nil, httpErrors.NewNotFoundError(errors.Wrap(err, "UserService.GetCurrentUser.FindById"))
	}

	// cache data to redis
	if err := u.redisRepo.Set(ctx, utils.GenerateRedisKey(basePrefix, id), cacheDuration, currentUser); err != nil {
		u.logger.WithError(err).Error("UserService.GetCurrentUser.redisRepo.Set")
	}

	return converter.ToUserResponse(currentUser), nil
}

// Update update current user
func (u *userService) Update(ctx context.Context, user *model.User) (*dto.UserResponse, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	updatedUser, err := u.pgRepo.Update(ctx, user)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	// delete redis cache data
	if err := u.redisRepo.Delete(ctx, utils.GetRedisKey(basePrefix, user.Id.String())); err != nil {
		u.logger.WithError(err).Error("UserService.Update.redisRepo.Delete")
	}

	return converter.ToUserResponse(updatedUser), nil
}

// UploadAvatar upload file
func (u *userService) UploadAvatar(ctx context.Context, file *model.UserUploadInput, id uuid.UUID) (*dto.UserResponse, error) {
	// upload to aws s3
	uploadInfo, err := u.awsRepo.PutObject(ctx, file)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	// generate url
	avatarURL := u.generateAWSMinioURL(file.BucketName, uploadInfo.Key)

	// update to database
	updatedUser, err := u.pgRepo.Update(ctx, &model.User{Id: id, Avatar: &avatarURL})

	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}
	return converter.ToUserResponse(updatedUser), nil
}

func (u *userService) generateAWSMinioURL(bucket string, key string) string {
	return fmt.Sprintf("%s/minio/%s/%s", u.cfg.AWS.Endpoint, bucket, key)
}
