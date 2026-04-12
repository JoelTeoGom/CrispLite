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

// ListConversations godoc
// @Summary      List conversations
// @Description  Returns all conversations for the authenticated user
// @Tags         chat
// @Produce      json
// @Security     BearerAuth
// @Success      200
// @Failure      401  {string}  string  "Unauthorized"
// @Router       /api/conversations [get]
func (h *ChatHandler) ListConversations(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

// LoadMessages godoc
// @Summary      Load messages
// @Description  Returns messages for a conversation
// @Tags         chat
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Conversation ID"
// @Success      200
// @Failure      401  {string}  string  "Unauthorized"
// @Router       /api/conversations/{id}/messages [get]
func (h *ChatHandler) LoadMessages(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
