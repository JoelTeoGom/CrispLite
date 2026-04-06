package outbound

import (
	"context"
	"crisplite/internal/domain"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	GetRefreshToken(ctx context.Context, hashedToken string) (*domain.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
}
