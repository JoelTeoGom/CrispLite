package rest

import (
	"crisplite/internal/adapter/inbound/rest/middleware"
	"crisplite/internal/port/outbound"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, uh *UserHandler, ch *ChatHandler, logger outbound.Logger, tokenService outbound.TokenService) http.Handler {
	mux.HandleFunc("POST /api/users", uh.CreateUser)
	mux.HandleFunc("POST /api/users/{id}/contacts", uh.AddContact)
	mux.HandleFunc("DELETE /api/users/{id}/contacts/{contactId}", uh.RemoveContact)

	mux.HandleFunc("GET /api/users/{id}/conversations", ch.ListConversations)
	mux.HandleFunc("GET /api/conversations/{id}/messages", ch.LoadMessages)

	var handler http.Handler = mux
	handler = middleware.Auth(tokenService, handler)
	handler = middleware.Logger(logger, handler)

	return handler
}
