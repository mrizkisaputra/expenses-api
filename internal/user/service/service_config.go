package service

import (
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/sirupsen/logrus"
)

// ServiceConfig will hold repositories that will eventually be injected into
// this service layer
type ServiceConfig struct {
	UserPostgresRepository user.UserPostgresRepository
	UserRedisRepository    user.UserRedisRepository
	AwsUserRepository      user.AWSUserRepository
	Logger                 *logrus.Logger
	Config                 *config.Config
}
