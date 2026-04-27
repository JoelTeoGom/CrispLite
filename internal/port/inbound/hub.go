package inbound

import (
	"context"

	"github.com/gorilla/websocket"
)

type Hub interface {
	Register(ctx context.Context, wsConn *websocket.Conn, userID string, device string) (string, error)
	Unregister(ctx context.Context, wsConn *websocket.Conn, userID string, connId string) error
}
