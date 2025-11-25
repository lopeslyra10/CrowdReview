package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from env vars.
type Config struct {
	AppPort           string
	DatabaseURL       string
	RedisURL          string
	JWTSecret         string
	RefreshSecret     string
	TokenTTL          time.Duration
	RefreshTTL        time.Duration
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

// LoadConfig loads environment variables and parses basic types.
func LoadConfig() Config {
	_ = godotenv.Load()

	tokenTTL := mustParseInt("TOKEN_TTL_MINUTES", 30)
	refreshTTL := mustParseInt("REFRESH_TTL_HOURS", 24)

	return Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/crowdreview?sslmode=disable"),
		RedisURL:          getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret"),
		RefreshSecret:     getEnv("REFRESH_SECRET", "dev-refresh-secret"),
		TokenTTL:          time.Duration(tokenTTL) * time.Minute,
		RefreshTTL:        time.Duration(refreshTTL) * time.Hour,
		RateLimitRequests: mustParseInt("RATE_LIMIT_REQUESTS", 20),
		RateLimitWindow:   time.Duration(mustParseInt("RATE_LIMIT_WINDOW", 60)) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustParseInt(key string, fallback int) int {
	v := getEnv(key, "")
	if v == "" {
		return fallback
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("invalid int for %s, using fallback %d", key, fallback)
		return fallback
	}
	return num
}
