package ws

import (
	"context"
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	onlineClients map[string][]*Connection
	chatService   inbound.ChatService
	mu            sync.RWMutex
	pubsub        outbound.PubSub
}

func NewHub(chatService inbound.ChatService, pubsub outbound.PubSub) *Hub {
	return &Hub{
		onlineClients: make(map[string][]*Connection),
		chatService:   chatService,
		pubsub:        pubsub,
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
	redisChannel := "user:" + userId
	h.pubsub.Subscribe(ctx, redisChannel)
	go conn.StartConnection(ctx)

	return connId, nil
}

func (h *Hub) Unregister(ctx context.Context, wsConn *websocket.Conn, userId, connId string) error {
	if wsConn == nil {
		return ErrNoWsConnProvided
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	conns, exists := h.onlineClients[userId]
	if !exists {
		return domain.ErrUserNotFound
	}

	for i, c := range conns {
		if c.ConnID == connId {
			closeMsg := websocket.FormatCloseMessage(
				websocket.CloseNormalClosure,
				"connection closed by client",
			)
			deadline := time.Now().Add(writeWait)
			_ = c.conn.WriteControl(websocket.CloseMessage, closeMsg, deadline)
			_ = c.conn.Close()
			h.onlineClients[userId] = append(conns[:i], conns[i+1:]...)
			break
		}
	}

	if len(h.onlineClients[userId]) == 0 {
		delete(h.onlineClients, userId)
		redisChannel := "user:" + userId
		h.pubsub.Unsubscribe(ctx, redisChannel)
	}
	return nil
}

func (h *Hub) Shutdown(ctx context.Context) {
	defer ctx.Done()
	h.mu.Lock()
	defer h.mu.Unlock()

	closeMsg := websocket.FormatCloseMessage(
		websocket.CloseGoingAway,
		"server shutting down",
	)
	deadline := time.Now().Add(writeWait)
	for userId, conns := range h.onlineClients {
		for _, c := range conns {
			_ = c.conn.WriteControl(websocket.CloseMessage, closeMsg, deadline)
			_ = c.conn.Close()
		}
		redisChannel := "user:" + userId
		h.pubsub.Unsubscribe(ctx, redisChannel)
		delete(h.onlineClients, userId)
	}
}
