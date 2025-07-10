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
var App Config

// LoadConfig reads configuration from environment variables.
func LoadConfig() (err error) {
	_, err = godotenv.Read(".env")
	if err != nil {
		fmt.Println("No .env file found, using environment variables")
	} else {
		err := godotenv.Load()
		if err != nil {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	App = Config{
		ServerPort:   os.Getenv("SERVER_PORT"),
		DBSource:     os.Getenv("DB_SOURCE"),
		RedisAddress: os.Getenv("REDIS_ADDRESS"),
		GinMode:      os.Getenv("GIN_MODE"),
		TrustedProxy: os.Getenv("TRUSTED_PROXY"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
	return nil
}
