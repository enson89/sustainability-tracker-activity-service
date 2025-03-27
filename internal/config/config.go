package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
}

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	if dbURL == "" || redisAddr == "" {
		return nil, errors.New("DATABASE_URL or REDIS_ADDR not set")
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	return &Config{
		DatabaseURL:   dbURL,
		RedisAddr:     redisAddr,
		RedisPassword: redisPassword,
	}, nil
}
