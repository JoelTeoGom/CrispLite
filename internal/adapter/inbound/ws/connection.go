package ws

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingInterval   = 30 * time.Second
	maxMessageSize = 512
)

type Connection struct {
	Client      *domain.Client
	chatService inbound.ChatService
}

// readPump reads messages from the websocket and forwards them to the chat service.
// It also handles pong messages to detect dead connections.
func (c *Connection) readPump() {
	// TODO: implement
	// - Set read limit (maxMessageSize)
	// - Set read deadline (pongWait)
	// - Set pong handler (reset read deadline on each pong)
	// - Loop: ReadMessage → unmarshal → chatService.Send
	// - On exit: hub.Disconnect + conn.Close
}

// writePump writes messages from the send channel to the websocket.
// It also sends periodic pings to detect dead connections.
func (c *Connection) writePump() {
	// TODO: implement
	// - Create ticker (pingInterval)
	// - Loop select:
	//   - case msg from c.send → set write deadline → conn.WriteMessage
	//   - case ticker → set write deadline → conn.WriteMessage(PingMessage)
	// - On exit: ticker.Stop + conn.Close
}
