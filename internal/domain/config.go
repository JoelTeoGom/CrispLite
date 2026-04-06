package domain

import "time"

type Env string

const (
	EnvLocal       Env = "local"
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
)

type Config struct {
	Env      Env
	Database DatabaseConfig
	Server   ServerConfig
	Batcher  BatcherConfig
}

type DatabaseConfig struct {
	URL          string
	MaxConns     int
	MinConns     int
	MaxConnLife  time.Duration
	MaxConnIdle  time.Duration
}

type ServerConfig struct {
	Port      string
	JWTSecret string
}

type BatcherConfig struct {
	Size     int
	Interval time.Duration
}
