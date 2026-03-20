package ws

import (
	"fmt"
	"net/http"
	"time"

	"crisplite/internal/application"
	"crisplite/internal/domain"
	"crisplite/internal/ports"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	repo ports.MessageRepository
}

func NewHandler(repo ports.MessageRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	messages := make(chan domain.Message, 10)
	go application.Batcher(messages, 10, 500*time.Millisecond, h.repo)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			close(messages)
			break
		}
		fmt.Printf("Received: %s\n", message)

		messages <- domain.Message{Content: string(message)}

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			close(messages)
			break
		}
	}
}
