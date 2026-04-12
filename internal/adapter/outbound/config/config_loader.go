package config

import (
	"bufio"
	"crisplite/internal/domain"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigLoader struct{}

func NewConfigLoader(path string) (*ConfigLoader, error) {
	if err := loadEnvFile(path); err != nil {
		return nil, err
	}
	return &ConfigLoader{}, nil
}

func (c *ConfigLoader) Load() (*domain.Config, error) {
	batchSize, _ := strconv.Atoi(getEnv("BATCH_SIZE", "100"))
	intervalMs, _ := strconv.Atoi(getEnv("BATCH_INTERVAL_MS", "200"))
	maxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNS", "20"))
	minConns, _ := strconv.Atoi(getEnv("DB_MIN_CONNS", "2"))
	maxConnLifeMin, _ := strconv.Atoi(getEnv("DB_MAX_CONN_LIFE_MIN", "30"))
	maxConnIdleMin, _ := strconv.Atoi(getEnv("DB_MAX_CONN_IDLE_MIN", "5"))

	return &domain.Config{
		Env: parseEnv(getEnv("APP_ENV", "local")),
		Database: domain.DatabaseConfig{
			URL:         getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/crisplite"),
			MaxConns:    maxConns,
			MinConns:    minConns,
			MaxConnLife: time.Duration(maxConnLifeMin) * time.Minute,
			MaxConnIdle: time.Duration(maxConnIdleMin) * time.Minute,
		},
		Server: domain.ServerConfig{
			Port:          getEnv("SERVER_PORT", "8080"),
			JWTSecret:     getEnv("JWT_SECRET", ""),
			AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		},
		Batcher: domain.BatcherConfig{
			Size:     batchSize,
			Interval: time.Duration(intervalMs) * time.Millisecond,
		},
		Redis: domain.RedisConfig{
			URI: getEnv("REDIS_URI", "redis://localhost:6379/1"),
		},
	}, nil
}

func parseEnv(value string) domain.Env {
	switch strings.ToLower(value) {
	case "development":
		return domain.EnvDevelopment
	case "production":
		return domain.EnvProduction
	default:
		return domain.EnvLocal
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	return scanner.Err()
}
