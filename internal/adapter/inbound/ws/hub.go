package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	onlineClients map[string][]*websocket.Conn
	mu            sync.RWMutex
}
