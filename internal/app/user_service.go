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
	if err := user.Validate(); err != nil {
		return "", err
	}
	if _, err := s.userRepo.CheckUserExists(ctx, user.Username); err != nil {
		return "", err
	}

	password, err := domain.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = password
	return s.userRepo.Save(ctx, user)
}

func (s *UserService) AddContact(ctx context.Context, userID, contactID string) error {
	if err := s.userRepo.AddContact(ctx, userID, contactID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveContact(ctx context.Context, userID, contactID string) error {
	if err := s.userRepo.RemoveContact(ctx, userID, contactID); err != nil {
		return err
	}
	return nil
}
