package rest

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userService inbound.UserService
}

func NewUserHandler(us inbound.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user := &domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := h.userService.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	var req struct {
		ContactID string `json:"contact_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.AddContact(userID, req.ContactID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	contactID := r.PathValue("contactId")

	if err := h.userService.RemoveContact(userID, contactID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
