package outbound

import (
	"context"
	"crisplite/internal/domain"
)

type TokenService interface {
	GenerateAccessToken(claims *domain.Claims) (string, error)
	GenerateRefreshToken(length int) string
	ValidateToken(token string) (*domain.Claims, error)
	AddClaimsToContext(ctx context.Context, claims *domain.Claims) context.Context
	ClaimsFromContext(ctx context.Context) (*domain.Claims, bool)
}
