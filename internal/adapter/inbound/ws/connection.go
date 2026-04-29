package ws

import (
	"context"
	"crisplite/internal/port/inbound"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingInterval   = 30 * time.Second
	maxMessageSize = 512
)

type Connection struct {
	ConnID      string
	UserID      string
	Device      string
	CreatedAt   time.Time
	conn        *websocket.Conn
	chatService inbound.ChatService
}

func (c *Connection) StartConnection(ctx context.Context) error {
	defer c.conn.Close()
	go c.writePump(ctx)

	err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		fmt.Println("Error setting write deadline:", err)

	}
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
	}

	return nil
}

func (c *Connection) readPump(ctx context.Context) {

	// ticker := time.NewTicker(pingInterval)
	// defer ticker.Stop()
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		return
	// 	case <- ticker.C:	//ping pong
	// 		pingMessage := []byte
	// 		writePump(ctx, pingMessage)

	// 	case <-conn.ReadMessage():
	// 		if
	// 	}

	for {

	}

}

// writePump sends periodic ping messages to the client to keep the connection alive and detect disconnections.
func (c *Connection) writePump(ctx context.Context) {
	pumpTicker := time.NewTicker(pingInterval)
	defer pumpTicker.Stop()

	for {

		// if err != nil {
		// 	fmt.Println("Error reading message:", err)
		// 	break
		// }

	}
}
