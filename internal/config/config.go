package config

import (
	"os"
)

type Config struct {
	Port            string
	DatabaseURL     string
	OrderServiceURL string
}

func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "3005"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		OrderServiceURL: getEnv("ORDER_SERVICE_URL", "http://localhost:3003"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
