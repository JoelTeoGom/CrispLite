package rest

import (
	"crisplite/internal/adapter/inbound/rest/middleware"
	"crisplite/internal/port/outbound"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, ah *AuthHandler, uh *UserHandler, ch *ChatHandler, logger outbound.Logger, tokenService outbound.TokenService) http.Handler {
	// Public routes (no auth)
	mux.HandleFunc("POST /api/auth/login", ah.Login)
	mux.HandleFunc("POST /api/auth/register", ah.Register)
	mux.HandleFunc("POST /api/auth/refresh", ah.RefreshToken)
	mux.HandleFunc("POST /api/auth/revoke", ah.RevokeToken)
	mux.HandleFunc("POST /api/auth/logout", ah.Logout)

	// Protected routes (auth required)
	protected := http.NewServeMux()
	protected.HandleFunc("POST /api/users/{id}/contacts", uh.AddContact)
	protected.HandleFunc("DELETE /api/users/{id}/contacts/{contactId}", uh.RemoveContact)
	protected.HandleFunc("GET /api/users/{id}/conversations", ch.ListConversations)
	protected.HandleFunc("GET /api/conversations/{id}/messages", ch.LoadMessages)

	mux.Handle("/api/", middleware.Auth(tokenService, protected))

	// Logger wraps everything
	return middleware.Logger(logger, mux)
}
