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
	JWTSecret    string
}

// AppConfig LoadConfig reads configuration from environment variables.
var AppConfig Config

// LoadConfig reads configuration from environment variables.
func LoadConfig() (config Config, err error) {
	_, err = godotenv.Read(".env")
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	} else {
		err := godotenv.Load()
		if err != nil {
			return config, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	config.ServerPort = os.Getenv("SERVER_PORT")
	config.DBSource = os.Getenv("DB_SOURCE")
	config.RedisAddress = os.Getenv("REDIS_ADDRESS")
	config.GinMode = os.Getenv("GIN_MODE")
	config.TrustedProxy = os.Getenv("TRUSTED_PROXY")
	config.JWTSecret = os.Getenv("JWT_SECRET")
	AppConfig = config // Set the global AppConfig
	return config, nil
}
