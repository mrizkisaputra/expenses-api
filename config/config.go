package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// App config
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Logger   LoggerConfig
	AWS      AwsConfig
}

// Server config
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Mode         string
	SSL          bool
	JWTSecretKey string
}

// Postgresql config
type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Dbname   string
}

// Logger config
type LoggerConfig struct {
	Level       string
	Caller      bool
	Encoding    string
	Development bool
}

type AwsConfig struct {
	Endpoint string
}

func NewAppConfig(configPath string) (*Config, error) {
	// get config path for docker, staging, production and default local development
	configFiles := map[string]string{
		"docker":     "./config/config-docker",
		"staging":    "",
		"production": "",
	}

	filename, isOk := configFiles[configPath]
	if !isOk {
		//return nil, fmt.Errorf("environment '%s' is not recognized", configPath)
		filename = "./config/config-local"
	}

	// load config file from given path
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file '%s' not found", filename)
		}
		return nil, fmt.Errorf("error reading config file, %v", err)
	}

	// parse config file
	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return cfg, nil
}
