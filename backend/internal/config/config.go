package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr           string
	DatabaseURL        string
	JWTSecret          string
	AccessTTL          time.Duration
	RefreshTTL         time.Duration
	CORSOrigins        []string
	UploadDir          string
	MaxUploadBytes     int64
	MigrationsPath     string
}

func Load() (Config, error) {
	accessMin, err := strconv.Atoi(getenv("JWT_ACCESS_TTL_MINUTES", "15"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_ACCESS_TTL_MINUTES: %w", err)
	}
	refreshHours, err := strconv.Atoi(getenv("JWT_REFRESH_TTL_HOURS", "168"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_REFRESH_TTL_HOURS: %w", err)
	}
	maxMB, err := strconv.Atoi(getenv("MAX_UPLOAD_MB", "5"))
	if err != nil {
		return Config{}, fmt.Errorf("MAX_UPLOAD_MB: %w", err)
	}

	secret := getenv("JWT_SECRET", "")
	if secret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	dbURL := getenv("DATABASE_URL", "")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	origins := strings.Split(getenv("CORS_ORIGINS", "http://localhost:4200"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return Config{
		HTTPAddr:       getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:    dbURL,
		JWTSecret:      secret,
		AccessTTL:      time.Duration(accessMin) * time.Minute,
		RefreshTTL:     time.Duration(refreshHours) * time.Hour,
		CORSOrigins:    origins,
		UploadDir:      getenv("UPLOAD_DIR", "./uploads"),
		MaxUploadBytes: int64(maxMB) * 1024 * 1024,
		MigrationsPath: getenv("MIGRATIONS_PATH", "file://migrations"),
	}, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
