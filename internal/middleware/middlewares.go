package middleware

import (
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/sirupsen/logrus"
)

type MiddlewareConfig struct {
	Logger *logrus.Logger
	Config *config.Config
}

// MiddlewareManager defines methods middleware
type MiddlewareManager struct {
	logger *logrus.Logger
	cfg    *config.Config
}

// NewMiddlewareManager is a factory function for instance MiddlewareManager
func NewMiddlewareManager(config *MiddlewareConfig) *MiddlewareManager {
	return &MiddlewareManager{
		logger: config.Logger,
		cfg:    config.Config,
	}
}
