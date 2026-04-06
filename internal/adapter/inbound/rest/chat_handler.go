package rest

import (
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"net/http"
)

type ChatHandler struct {
	chatService inbound.ChatService
	logger      outbound.Logger
}

func NewChatHandler(cs inbound.ChatService, logger outbound.Logger) *ChatHandler {
	return &ChatHandler{chatService: cs, logger: logger}
}

func (h *ChatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (h *ChatHandler) LoadMessages(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
