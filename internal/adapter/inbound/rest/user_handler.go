package rest

import (
	"crisplite/internal/port/inbound"
	"net/http"
)

type UserHandler struct {
	userService inbound.UserService
}

func NewUserHandler(us inbound.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (h *UserHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}

func (h *UserHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	// TODO implement
}
