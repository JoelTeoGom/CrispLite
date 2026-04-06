package outbound

import "context"

type Claims struct {
	UserID string
	Role   string
}

type TokenService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(length int) string
	ValidateToken(token string) (*Claims, error)
	AddClaimsToContext(ctx context.Context, claims *Claims) context.Context
	ClaimsFromContext(ctx context.Context) (*Claims, bool)
}
