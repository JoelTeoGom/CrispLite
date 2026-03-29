package ws

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub         *Hub
	chatService inbound.ChatService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	// TODO: extract userId from request (auth JWT)
	userId := r.URL.Query().Get("userId")
	uuid := uuid.New().String()
	createdAt := time.Now()

	client := &domain.Client{
		ConnID:    uuid,
		UserID:    userId,
		Device:    "web", // TODO: determine device type
		Conn:      wsConn,
		CreatedAt: createdAt,
	}

	conn := NewConnection(h.hub, wsConn, userId)
	h.hub.Connect(userId, client)

	go conn.readPump()
	go conn.writePump()
}
