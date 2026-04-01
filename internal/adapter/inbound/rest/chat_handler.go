package rest

import (
	"crisplite/internal/port/inbound"
	"net/http"
)

type ChatHandler struct {
	chatService inbound.ChatService
}

func NewChatHandler(cs inbound.ChatService) *ChatHandler {
	return &ChatHandler{chatService: cs}
}

func (h *ChatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (h *ChatHandler) LoadMessages(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
