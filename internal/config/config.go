package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	MongoURI    string
	Database    string
	JWTSecret   string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	env := getEnv("ENV", "development")
	mongoURI := getEnv("MONGO_URI", "mongodb://mongo:27017")
	dbName := getEnv("MONGO_DB", "astro")
	jwtSecret := getEnv("JWT_SECRET", "supersecretkey")

	return &Config{
		Port:        port,
		Environment: env,
		MongoURI:    mongoURI,
		Database:    dbName,
		JWTSecret:   jwtSecret,
	}
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
