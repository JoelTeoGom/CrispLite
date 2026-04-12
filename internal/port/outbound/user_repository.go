package outbound

import (
	"context"
	"crisplite/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (string, error)
	AddContact(ctx context.Context, userID, contactID string) error
	RemoveContact(ctx context.Context, userID, contactID string) error
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	SearchUsers(ctx context.Context, query string, limit, offset int) ([]domain.UserSummary, error)
}
