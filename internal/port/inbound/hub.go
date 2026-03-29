package inbound

import "crisplite/internal/domain"

type Hub interface {
	Connect(userId string, client *domain.Client) error
	Disconnect(userId string, connID string) error
}
