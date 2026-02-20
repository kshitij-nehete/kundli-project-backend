package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")

	return &Config{
		Port:        port,
		Environment: env,
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
