package app

import (
	"crisplite/internal/adapter/inbound/ws"
	"crisplite/internal/domain"
	"crisplite/internal/port/outbound"
)

type ChatService struct {
	Hub     *ws.Hub
	Batcher Batcher
	logger  outbound.Logger
}

func NewChatService(messageRepo outbound.MessageRepository, batcher Batcher, logger outbound.Logger) *ChatService {
	return &ChatService{
		Hub:     ws.NewHub(),
		Batcher: batcher,
		logger:  logger,
	}
}

func (s *ChatService) Send(msg *domain.Message) error {
	// TODO: implement
	//take te message and send it to batcher and redis

	return nil
}

func (s *ChatService) Deliver(msg *domain.Message) error {
	// TODO: implement

	// take HUB and send messages (we can do logic on them)
	return nil
}

// SEND:
//   Connection.readPump → Handler → ChatService.Send(msg)
//                                       ├──→ msg channel → Batcher → Postgres
//                                       └──→ Broker.Publish(msg) → Redis

// RECEIVE:
//   Redis → Worker/Subscriber → ChatService.Deliver(msg)
//                                    └──→ Hub.SendToUser(userId, data)
//                                             ├──→ conn1.send chan
//                                             └──→ conn2.send chan
