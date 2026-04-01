package inbound

import "crisplite/internal/domain"

type UserService interface {
	CreateUser(user *domain.User) error
	AddContact(userID, contactID string) error
	RemoveContact(userID, contactID string) error
}
