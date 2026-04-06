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

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest   true  "Register request"
// @Success      200   {object}  RegisterResponse
// @Failure      400   {string}  string  "invalid request body"
// @Failure      500   {string}  string  "failed to register user"
// @Router       /api/auth/register [post]
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

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest   true  "Login request"
// @Success      200   {object}  RegisterResponse
// @Failure      400   {string}  string  "invalid request body"
// @Failure      500   {string}  string  "failed to login user"
// @Router       /api/auth/login [post]
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

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Returns new access and refresh tokens given a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      object{refresh_token=string}  true  "Refresh token request"
// @Success      200   {object}  RegisterResponse
// @Failure      400   {string}  string  "invalid request body"
// @Failure      500   {string}  string  "failed to refresh token"
// @Router       /api/auth/refresh [post]
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

// Logout godoc
// @Summary      Logout user
// @Description  Revokes the refresh token to log the user out
// @Tags         auth
// @Accept       json
// @Param        body  body      object{refresh_token=string}  true  "Logout request"
// @Success      204
// @Failure      400   {string}  string  "invalid request body"
// @Failure      500   {string}  string  "failed to revoke token"
// @Router       /api/auth/logout [post]
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

// RevokeToken godoc
// @Summary      Revoke a refresh token
// @Description  Marks a refresh token as revoked
// @Tags         auth
// @Accept       json
// @Param        body  body      object{refresh_token=string}  true  "Revoke token request"
// @Success      204
// @Failure      400   {string}  string  "invalid request body"
// @Failure      500   {string}  string  "failed to revoke token"
// @Router       /api/auth/revoke [post]
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
