package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv    string
	AppPort   string
	JWTSecret string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	AIProvider string
}

func Load() Config {
	return Config{
		AppEnv:    env("APP_ENV", "development"),
		AppPort:   env("APP_PORT", "8080"),
		JWTSecret: env("JWT_SECRET", "local-development-secret"),

		DBHost:     env("DB_HOST", "localhost"),
		DBPort:     env("DB_PORT", "3306"),
		DBUser:     env("DB_USER", "ai_nutrition"),
		DBPassword: env("DB_PASSWORD", "ai_nutrition"),
		DBName:     env("DB_NAME", "ai_nutrition"),

		RedisAddr:     env("REDIS_ADDR", "localhost:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
		RedisDB:       envInt("REDIS_DB", 0),

		AIProvider: env("AI_PROVIDER", "mock"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
