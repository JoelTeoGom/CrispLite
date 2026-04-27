package ws

import (
	"context"
	"crisplite/internal/port/inbound"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	onlineClients map[string][]*Connection
	chatService   inbound.ChatService
	mu            sync.RWMutex
}

func NewHub(chatService inbound.ChatService) *Hub {
	return &Hub{
		onlineClients: make(map[string][]*Connection),
		chatService:   chatService,
	}
}

func (h *Hub) Register(ctx context.Context, wsConn *websocket.Conn, userId, device string) (string, error) {
	if wsConn == nil {
		return "", ErrNoWsConnProvided
	}
	connId := uuid.New().String()
	conn := &Connection{
		ConnID:      connId,
		UserID:      userId,
		Device:      device,
		CreatedAt:   time.Now(),
		conn:        wsConn,
		chatService: h.chatService,
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.onlineClients[userId] = append(h.onlineClients[userId], conn)

	go conn.StartConnection(ctx)

	return connId, nil
}

func (h *Hub) Unregister(ctx context.Context, wsConn *websocket.Conn, userId, connId string) error {
	if wsConn == nil {
		return ErrNoWsConnProvided
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
