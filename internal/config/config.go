package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Port          string
	DatabasePath  string
	UploadDir     string
	PublicSiteURL string
	JWTSecret     string
	AdminUsername string
	AdminPassword string
}

func Load() Config {
	return Config{
		Port:          env("PORT", "8080"),
		DatabasePath:  env("DATABASE_PATH", filepath.Join("data", "app.db")),
		UploadDir:     env("UPLOAD_DIR", "public_uploads"),
		PublicSiteURL: env("PUBLIC_SITE_URL", ""),
		JWTSecret:     env("JWT_SECRET", "change-this-secret-in-production"),
		AdminUsername: env("ADMIN_USERNAME", "admin"),
		AdminPassword: env("ADMIN_PASSWORD", "admin123456"),
	}
}

func (c Config) Address() string {
	return ":" + c.Port
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func EnvInt(key string, fallback int) int {
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
