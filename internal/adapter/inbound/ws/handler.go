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
	chatService inbound.ChatService // puerto inbound, no el canal directamente
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
				fmt.Println("conexión cerrada limpiamente")
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

		//we send that to the general message buffer
		h.messages <- m
		fmt.Printf("Whatsapp enviado por %s {contenido: %s} a %s\n", m.SenderId, m.Content, m.ReceiverId)
	}
}

// func main() {
// 	http.HandleFunc("/ws", wsHandler)
// 	fmt.Println("WebSocket server started on :8080")
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Println("Error starting server:", err)
// 	}
// }
