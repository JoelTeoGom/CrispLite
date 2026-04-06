package rest

import (
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userService inbound.UserService
	logger      outbound.Logger
}

type CreateUserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewUserHandler(us inbound.UserService, logger outbound.Logger) *UserHandler {
	return &UserHandler{userService: us, logger: logger}
}

func (h *UserHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ContactID string `json:"contact_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.userService.AddContact(r.Context(), req.ContactID)
	if err != nil {
		http.Error(w, "failed to add contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    string `json:"user_id"`
		ContactID string `json:"contact_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.userService.RemoveContact(r.Context(), req.ContactID)
	if err != nil {
		http.Error(w, "failed to remove contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
