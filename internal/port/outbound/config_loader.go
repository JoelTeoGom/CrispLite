package outbound

import "crisplite/internal/domain"

type ConfigLoader interface {
	Load() (*domain.Config, error)
}
