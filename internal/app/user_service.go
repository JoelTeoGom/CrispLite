package app

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
)

type UserService struct {
	userRepo    outbound.UserRepository
	authService outbound.TokenService
	logger      outbound.Logger
}

func NewUserService(userRepo outbound.UserRepository, authService outbound.TokenService, logger outbound.Logger) *UserService {
	return &UserService{userRepo: userRepo, authService: authService, logger: logger}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", err
	}

	password, err := domain.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = password

	return s.userRepo.Save(ctx, user)
}

func (s *UserService) Login(ctx context.Context, username, password string) (*domain.RegisterResponse, error) {

	return nil, nil
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.RegisterResponse, error) {

	return nil, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error) {

	return nil, nil
}

func (s *UserService) RevokeToken(ctx context.Context, refreshToken string) error {

	return nil
}

func (s *UserService) AddContact(ctx context.Context, contactID string) error {
	claims, ok := s.authService.ClaimsFromContext(ctx)
	if !ok {
		return domain.ErrUnauthorized
	}
	userID := claims.UserID
	if err := s.userRepo.AddContact(ctx, userID, contactID); err != nil {
		return err
	}
	return nil
}

func (s *UserService) RemoveContact(ctx context.Context, contactID string) error {
	claims, ok := s.authService.ClaimsFromContext(ctx)
	if !ok {
		return domain.ErrUnauthorized
	}
	userID := claims.UserID
	if err := s.userRepo.RemoveContact(ctx, userID, contactID); err != nil {
		return err
	}
	return nil
}
