package app

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
	"time"
)

type UserService struct {
	userRepo    outbound.UserRepository
	authRepo    outbound.AuthRepository
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
	if username == "" || password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	password, err = domain.HashPassword(password)
	if err != nil {
		return nil, err
	}
	if password != user.Password {
		return nil, domain.ErrInvalidCredentials
	}

	claims := &domain.Claims{
		UserID: user.ID,
		Role:   "user",
	}
	accessToken, err := s.authService.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken := s.authService.GenerateRefreshToken(32)
	hashedRefreshToken, err := domain.HashPassword(refreshToken)

	domainRefreshToken := &domain.RefreshToken{
		UserID:      user.ID,
		HashedToken: hashedRefreshToken,
	}

	err = s.authRepo.SaveRefreshToken(ctx, domainRefreshToken)
	if err != nil {
		return nil, err
	}
	return &domain.RegisterResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RegisterUser(ctx context.Context, user *domain.User) (*domain.RegisterResponse, error) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	password, err := domain.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = password

	userId, err := s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	claims := &domain.Claims{
		UserID: userId,
		Role:   "user",
	}
	accessToken, err := s.authService.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken := s.authService.GenerateRefreshToken(32)
	hashedRefreshToken, err := domain.HashPassword(refreshToken)
	if err != nil {
		return nil, err
	}
	domainRefreshToken := &domain.RefreshToken{
		UserID:      userId,
		HashedToken: hashedRefreshToken,
	}

	err = s.authRepo.SaveRefreshToken(ctx, domainRefreshToken)
	if err != nil {
		return nil, err
	}
	return &domain.RegisterResponse{
		UserID:       userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshResponse, error) {
	if refreshToken == "" {
		return nil, domain.ErrInvalidToken
	}

	hashedToken, err := domain.HashPassword(refreshToken)
	if err != nil {
		return nil, err
	}

	storedToken, err := s.authRepo.GetRefreshToken(ctx, hashedToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if storedToken.Revoked {
		return nil, domain.ErrRevokedToken
	}

	if time.Since(storedToken.CreatedAt) > 7*24*time.Hour {
		return nil, domain.ErrExpiredToken
	}

	claims := &domain.Claims{
		UserID: storedToken.UserID,
		Role:   "user",
	}
	accessToken, err := s.authService.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken = s.authService.GenerateRefreshToken(32)
	hashedRefreshToken, err := domain.HashPassword(refreshToken)
	if err != nil {
		return nil, err
	}

	domainRefreshToken := &domain.RefreshToken{
		UserID:      storedToken.UserID,
		HashedToken: hashedRefreshToken,
	}

	err = s.authRepo.SaveRefreshToken(ctx, domainRefreshToken)
	if err != nil {
		return nil, err
	}
	return &domain.RefreshResponse{
		UserID:       storedToken.UserID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (s *UserService) RevokeToken(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return domain.ErrInvalidToken
	}

	hashedToken, err := domain.HashPassword(refreshToken)
	if err != nil {
		return err
	}

	err = s.authRepo.RevokeRefreshToken(ctx, hashedToken)
	if err != nil {
		return err
	}
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
