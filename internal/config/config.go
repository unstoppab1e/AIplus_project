package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
}

func NewConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/employees?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
