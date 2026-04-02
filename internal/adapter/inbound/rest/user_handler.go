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

type CreateUserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewUserHandler(us inbound.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var domainUser = &domain.User{
		Username: req.UserName,
		Password: req.Password,
	}
	userID, err := h.userService.CreateUser(r.Context(), domainUser)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	payload := map[string]string{"id": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)

}

func (h *UserHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (h *UserHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
