package postgres

import (
	"fmt"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

const (
	maxOpenConns    = 60  // Maksimal koneksi terbuka
	connMaxLifetime = 120 // Maksimal waktu hidup koneksi dalam detik
	maxIdleConns    = 30  // Maksimal koneksi idle
	connMaxIdleTime = 20  // Maksimal waktu idle koneksi dalam detik
)

// return new postgresql db instance
func NewPostgresConn(cfg *config.Config) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Postgres.Host,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Dbname,
		cfg.Postgres.Port,
	)
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger:                 logger.Default.LogMode(getLoggerLevel(cfg)),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "NewPostgresConn.Open")
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "NewPostgresConn.DB")
	}
	sqlDb.SetMaxOpenConns(maxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Second * connMaxLifetime)
	sqlDb.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
	sqlDb.SetMaxIdleConns(maxIdleConns)

	return db, nil
}

var fieldsLogrusLevelMap = map[string]logger.LogLevel{
	"info":   logger.Info,
	"warn":   logger.Warn,
	"error":  logger.Error,
	"silent": logger.Silent,
}

func getLoggerLevel(cfg *config.Config) logger.LogLevel {
	level, exist := fieldsLogrusLevelMap[cfg.Logger.Level]
	if !exist {
		level = logger.Info
	}
	return level
}
