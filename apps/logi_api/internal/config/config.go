package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config stores all configuration of the application.
// The values are read from environment variables.
type Config struct {
	ServerPort   string
	DBSource     string
	RedisAddress string
	GinMode      string
	TrustedProxy string
}

// LoadConfig reads configuration from environment variables.
func LoadConfig() (config Config, err error) {
	if err := godotenv.Load(); err != nil {
		return config, fmt.Errorf("error loading .env file: %w", err)
	}

	config.ServerPort = os.Getenv("SERVER_PORT")
	config.DBSource = os.Getenv("DB_SOURCE")
	config.RedisAddress = os.Getenv("REDIS_ADDRESS")
	config.GinMode = os.Getenv("GIN_MODE")
	config.TrustedProxy = os.Getenv("TRUSTED_PROXY")
	return config, nil
}
