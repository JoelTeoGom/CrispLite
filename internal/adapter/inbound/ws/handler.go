package ws

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println("connection closed cleanly")
			} else {
				log.Println("read error:", err)
			}
			break
		}
		var m domain.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		// send message to chat service
		err = h.chatService.Send(&m)
		if err != nil {
			log.Println("error scheduling message")
			continue
		}
		fmt.Printf("message sent by %s {content: %s} to %s\n", m.SenderId, m.Content, m.ReceiverId)
	}
}

