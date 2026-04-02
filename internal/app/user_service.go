package app

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
)

type UserService struct {
	userRepo outbound.UserRepository
}

func NewUserService(userRepo outbound.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	if user.Username == "" || user.Password == "" {
		return "", nil // TODO return proper error
	}
	if _, err := s.userRepo.CheckUserExists(ctx, user.Username); err != nil {
		return "", err
	}

	password := hashPassword(user.Password) // TODO implement password hashing
	user.Password = password
	return s.userRepo.Save(ctx, user)
}

func (s *UserService) AddContact(ctx context.Context, userID, contactID string) error {
	// TODO implement
	return nil
}

func (s *UserService) RemoveContact(ctx context.Context, userID, contactID string) error {
	// TODO implement
	return nil
}
