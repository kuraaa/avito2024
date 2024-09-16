// internal/config/config.go
package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerAddress string
	PostgresConn  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		PostgresConn:  os.Getenv("POSTGRES_CONN"),
	}, nil
}