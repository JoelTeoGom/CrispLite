package domain

import "errors"

var (
	ErrUsernameEmpty    = errors.New("username is required")
	ErrUsernameTooShort = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong  = errors.New("username must be at most 50 characters")
	ErrPasswordEmpty    = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong      = errors.New("password must be at most 72 characters")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserNotFound         = errors.New("user not found")
	ErrContactAlreadyExists = errors.New("contact already exists")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidToken         = errors.New("invalid token")
	ErrRevokedToken         = errors.New("token has been revoked")
	ErrExpiredToken         = errors.New("token has expired")
)
