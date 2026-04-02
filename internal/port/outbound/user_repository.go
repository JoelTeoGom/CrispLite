package outbound

import (
	"context"
	"crisplite/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (string, error)
	AddContact(ctx context.Context, userID, contactID string) error
	RemoveContact(ctx context.Context, userID, contactID string) error
	CheckUserExists(ctx context.Context, username string) (bool, error)
}
