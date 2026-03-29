package ws

import (
	"crisplite/internal/domain"
	"sync"
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

func (h *Hub) Connect(userId string, client *domain.Client) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.onlineClients[userId] = append(h.onlineClients[userId], client)

	return nil
}

func (h *Hub) Disconnect(userId string, connID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients := h.onlineClients[userId]
	for i, c := range clients {
		if c.ConnID == connID {
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
