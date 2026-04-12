package rest

import (
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userService  inbound.UserService
	tokenService outbound.TokenService
	logger       outbound.Logger
}

type CreateUserRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func NewUserHandler(us inbound.UserService, ts outbound.TokenService, logger outbound.Logger) *UserHandler {
	return &UserHandler{userService: us, tokenService: ts, logger: logger}
}

// AddContact godoc
// @Summary      Add a contact
// @Description  Adds a contact to the authenticated user
// @Tags         contacts
// @Accept       json
// @Security     BearerAuth
// @Param        body  body      object{contact_id=string}     true  "Contact to add"
// @Success      201
// @Failure      400   {string}  string  "invalid request body"
// @Failure      401   {string}  string  "Unauthorized"
// @Failure      500   {string}  string  "failed to add contact"
// @Router       /api/contacts [post]
func (h *UserHandler) AddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := h.tokenService.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	err := h.userService.AddContact(r.Context(), claims.UserID, req.ContactID)
	if err != nil {
		http.Error(w, "failed to add contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// RemoveContact godoc
// @Summary      Remove a contact
// @Description  Removes a contact from the authenticated user
// @Tags         contacts
// @Security     BearerAuth
// @Param        contactId  path      string  true  "Contact ID"
// @Success      204
// @Failure      401   {string}  string  "Unauthorized"
// @Failure      500   {string}  string  "failed to remove contact"
// @Router       /api/contacts/{contactId} [delete]
func (h *UserHandler) RemoveContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := h.tokenService.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	contactID := r.PathValue("contactId")
	if contactID == "" {
		http.Error(w, "missing contact id", http.StatusBadRequest)
		return
	}

	err := h.userService.RemoveContact(r.Context(), claims.UserID, contactID)
	if err != nil {
		http.Error(w, "failed to remove contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
