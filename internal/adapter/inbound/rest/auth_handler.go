package rest

import (
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
)

type authHandler struct {
	userService inbound.UserService
	logger      outbound.Logger
}
