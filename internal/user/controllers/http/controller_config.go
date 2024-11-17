package http

import (
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/sirupsen/logrus"
)

type ControllerConfig struct {
	Config      *config.Config
	Logger      *logrus.Logger
	AuthService user.AuthService
	UserService user.UserService
}
