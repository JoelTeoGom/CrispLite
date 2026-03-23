package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL   string
	ServerPort    string
	BatchSize     int
	BatchInterval time.Duration
}

func Load(path string) (*Config, error) {
	if err := loadEnvFile(path); err != nil {
		return nil, err
	}

	batchSize, _ := strconv.Atoi(getEnv("BATCH_SIZE", "100"))
	intervalMs, _ := strconv.Atoi(getEnv("BATCH_INTERVAL_MS", "200"))

	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/crisplite"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		BatchSize:     batchSize,
		BatchInterval: time.Duration(intervalMs) * time.Millisecond,
	}, nil
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
