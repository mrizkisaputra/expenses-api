package redis

import (
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		DB:           cfg.Redis.DB,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  cfg.Redis.PoolTimeout * time.Second,
		Password:     cfg.Postgres.Password,
	})

	return client
}
