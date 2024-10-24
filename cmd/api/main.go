package main

import (
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/server"
	"github.com/mrizkisaputra/expenses-api/pkg/db/postgres"
	"github.com/mrizkisaputra/expenses-api/pkg/logger"
	"log"
	"os"
)

func main() {

	// Initialize app config
	cfg, err := config.NewAppConfig(os.Getenv("config"))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Logger
	apiLogger := logger.NewLogrusLogger(cfg)

	// Intialize postgreSQL connection
	psqlDB, err := postgres.NewPostgresConn(cfg)
	if err != nil {
		apiLogger.Fatalf("Postgresql initialize: %v", err)
	}
	apiLogger.Info("PostgreSQL connected")

	// Initialize server
	s := server.NewServer(apiLogger, cfg, psqlDB)
	if err := s.Run(); err != nil {
		apiLogger.Fatal(err)
	}

}
