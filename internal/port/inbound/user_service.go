package inbound

import (
	"context"
	"crisplite/internal/domain"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (*domain.RegisterResponse, error)
	RegisterUser(ctx context.Context, user *domain.User) (*domain.RegisterResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error)
	RevokeToken(ctx context.Context, refreshToken string) error
	AddContact(ctx context.Context, userID, contactID string) error
	RemoveContact(ctx context.Context, userID, contactID string) error
}
