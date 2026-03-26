package inbound

import "crisplite/internal/domain"

type ChatService interface {
	Send(msg *domain.Message) error
	Deliver(msg *domain.Message) error
}
