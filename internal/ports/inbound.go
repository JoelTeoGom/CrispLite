package ports

import "crisplite/internal/domain"

type MessageService interface {
	HandleMessage(msg domain.Message) error
}
