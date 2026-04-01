package rest

import "net/http"

func RegisterRoutes(mux *http.ServeMux, uh *UserHandler, ch *ChatHandler) {
	mux.HandleFunc("POST /api/users", uh.userService.CreateUser())
	mux.HandleFunc("POST /api/users/{id}/contacts", uh.userService.AddContact)
	mux.HandleFunc("DELETE /api/users/{id}/contacts/{contactId}", uh.userService.RemoveContact)

	mux.HandleFunc("GET /api/users/{id}/conversations", ch.ListConversations)
	mux.HandleFunc("GET /api/conversations/{id}/messages", ch.LoadMessages)
}
