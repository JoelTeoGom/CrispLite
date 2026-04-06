package rest

import (
	"crisplite/internal/domain"
	"crisplite/internal/port/inbound"
	"crisplite/internal/port/outbound"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	userService inbound.UserService
	logger      outbound.Logger
}

func NewAuthHandler(us inbound.UserService, logger outbound.Logger) *AuthHandler {
	return &AuthHandler{userService: us, logger: logger}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	UserId       string `json:"user_id"`
	JwtToken     string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var domainUser = &domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	registerResponse, err := h.userService.RegisterUser(r.Context(), domainUser)
	if err != nil {
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	payload := &RegisterResponse{
		UserId:       registerResponse.UserID,
		JwtToken:     registerResponse.AccessToken,
		RefreshToken: registerResponse.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var domainUser = &domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	registerResponse, err := h.userService.RegisterUser(r.Context(), domainUser)
	if err != nil {
		http.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}

	payload := &RegisterResponse{
		UserId:       registerResponse.UserID,
		JwtToken:     registerResponse.AccessToken,
		RefreshToken: registerResponse.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	refreshResponse, err := h.userService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "failed to refresh token", http.StatusInternalServerError)
		return
	}

	payload := &RegisterResponse{
		UserId:       refreshResponse.UserID,
		JwtToken:     refreshResponse.AccessToken,
		RefreshToken: refreshResponse.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.userService.RevokeToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "failed to revoke token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := h.userService.RevokeToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "failed to revoke token", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
