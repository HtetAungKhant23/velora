package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	StorageBaseDir string
	StorageBaseURL string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		Port:           getEnv("PORT", "3003"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiry:      getDurationEnv("JWT_EXPIRY_HOURS", 24) * time.Hour,
		StorageBaseDir: getEnv("STORAGE_BASE_DIR", ""),
		StorageBaseURL: getEnv("STORAGE_BASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fallback
	}
	return time.Duration(n)
}
