package postgres

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool   *pgxpool.Pool
	logger outbound.Logger
}

func NewAuthRepo(pool *pgxpool.Pool, logger outbound.Logger) *AuthRepo {
	return &AuthRepo{pool: pool, logger: logger}
}

func (r *AuthRepo) SaveRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO tokens (hashed_token, user_id) VALUES ($1, $2)",
		token.HashedToken, token.UserID,
	)
	return err
}

func (r *AuthRepo) GetRefreshToken(ctx context.Context, hashedToken string) (*domain.RefreshToken, error) {
	var t domain.RefreshToken
	err := r.pool.QueryRow(ctx,
		"SELECT id, hashed_token, user_id, revoked, created_at FROM tokens WHERE hashed_token = $1",
		hashedToken,
	).Scan(&t.ID, &t.HashedToken, &t.UserID, &t.Revoked, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *AuthRepo) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE tokens SET revoked = true WHERE id = $1",
		tokenID,
	)
	return err
}
