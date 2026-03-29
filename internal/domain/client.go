package domain

import (
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ConnID    string // UUID generado server-side
	UserID    string // del JWT
	Device    string // "web", "app-ios", "app-android"
	Conn      *websocket.Conn
	CreatedAt time.Time
}
