package ws

import (
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	hub          inbound.Hub
	tokenService outbound.TokenService
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
	defer wsConn.Close()

	claims, ok := h.tokenService.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	connId, err := h.hub.Register(r.Context(), wsConn, claims.UserID, "web")
	if err != nil {
		http.Error(w, "Error registering client", http.StatusInternalServerError)
		return
	}
	defer h.hub.Unregister(r.Context(), wsConn, claims.UserID, connId)
}
