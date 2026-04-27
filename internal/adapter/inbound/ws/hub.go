package ws

import (
	"context"
	"crisplite/internal/domain"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	onlineClients map[string][]*domain.Client
	mu            sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		onlineClients: make(map[string][]*domain.Client),
	}
}

func (h *Hub) Register(ctx context.Context, wsConn *websocket.Conn, userId, device string) (string, error) {
	if wsConn == nil {
		return "", domain.ErrNoWsConnProvided
	}
	connId := uuid.New().String()
	createdAt := time.Now()
	client := &domain.Client{
		ConnID:    connId,
		UserID:    userId,
		Device:    device,
		Conn:      wsConn,
		CreatedAt: createdAt,
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onlineClients[userId] = append(h.onlineClients[userId], client)
	conn := &Connection{
		Client: client,
	}

	go conn.StartConnection(ctx)

	return connId, nil
}

func (h *Hub) Unregister(ctx context.Context, wsConn *websocket.Conn, userId, connId string) error {
	if wsConn == nil {
		return domain.ErrNoWsConnProvided
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.onlineClients[userId]
	for i, c := range clients {
		if c.ConnID == connId {
			h.onlineClients[userId] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	if len(h.onlineClients[userId]) == 0 {
		delete(h.onlineClients, userId)
		//TODO  the client is offline, we need to delete REDIS channel and notify other clients
	}

	return nil
}
