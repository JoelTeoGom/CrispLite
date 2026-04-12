package rest

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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
		if errors.Is(err, domain.ErrContactAlreadyExists) {
			http.Error(w, "contact already exists", http.StatusConflict)
			return
		}
		if errors.Is(err, domain.ErrInvalidContact) {
			http.Error(w, "invalid contact", http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to add contact", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// SearchUsers godoc
// @Summary      Search users by username
// @Description  Returns users matching the query (min 3 chars). Excludes the authenticated user.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        q       query     string  true   "Search text (min 3 chars)"
// @Param        limit   query     int     false  "Max results (default 20)"
// @Param        offset  query     int     false  "Offset (default 0)"
// @Success      200     {array}   domain.UserSummary
// @Failure      400     {string}  string  "search query too short"
// @Failure      401     {string}  string  "Unauthorized"
// @Failure      500     {string}  string  "failed to search users"
// @Router       /api/users/search [get]
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.tokenService.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	users, err := h.userService.SearchUsers(r.Context(), claims.UserID, query, limit, offset)
	if err != nil {
		http.Error(w, "failed to search users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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
