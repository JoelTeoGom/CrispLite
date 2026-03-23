package domain

import "time"

type Config struct {
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
	Port string
}

type BatcherConfig struct {
	Size     int
	Interval time.Duration
}
