package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/server"
	"github.com/mrizkisaputra/expenses-api/pkg/db/aws"
	"github.com/mrizkisaputra/expenses-api/pkg/db/postgres"
	"github.com/mrizkisaputra/expenses-api/pkg/db/redis"
	"github.com/mrizkisaputra/expenses-api/pkg/logger"
	"log"
	"os"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)

	// -----------------------------------------------------------------------------------------------------------
	// initializing app config
	cfg, err := config.NewAppConfig(os.Getenv("config"))
	if err != nil {
		log.Fatal(err)
	}

	// -----------------------------------------------------------------------------------------------------------
	// initializing Logger
	apiLogger := logger.NewLogrusLogger(cfg)

	// -----------------------------------------------------------------------------------------------------------
	// initializing postgreSQL connection
	psqlDB, err := postgres.NewPostgresConn(cfg)
	if err != nil {
		apiLogger.Fatalf("Postgresql initialize: %v", err)
	}
	apiLogger.Info("PostgreSQL connected")

	// -----------------------------------------------------------------------------------------------------------
	// initializing redis client
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		apiLogger.Fatal("Redis ping: %v", err)
	}
	apiLogger.Info("Redis connected")

	// -----------------------------------------------------------------------------------------------------------
	// initializing aws minio client
	awsClient, err := aws.NewAWSClient(cfg)
	if err != nil {
		apiLogger.Fatalf("Minio initializing: %v", err)
	}
	apiLogger.Info("Minio connected")

	// -----------------------------------------------------------------------------------------------------------
	// instance gin framework
	app := gin.New()

	// -----------------------------------------------------------------------------------------------------------
	// create a new instance server
	s := server.NewServer(&server.ServerConfig{
		App:         app,
		Cfg:         cfg,
		Logger:      apiLogger,
		Db:          psqlDB,
		RedisClient: redisClient,
		AwsClient:   awsClient,
	})

	if err := s.Run(); err != nil {
		apiLogger.Fatal(err)
	}

}
