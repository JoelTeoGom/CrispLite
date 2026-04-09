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
	secure      bool
}

func NewAuthHandler(us inbound.UserService, logger outbound.Logger, env domain.Env) *AuthHandler {
	return &AuthHandler{
		userService: us,
		logger:      logger,
		secure:      env == domain.EnvProduction,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	UserId   string `json:"user_id"`
	JwtToken string `json:"jwt_token"`
}

func (h *AuthHandler) setRefreshTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth",
		MaxAge:   7 * 24 * 60 * 60,
	})
}

func (h *AuthHandler) clearRefreshTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth",
		MaxAge:   -1,
	})
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user and returns access token. Refresh token is set as httpOnly cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Register request"
// @Success      200   {object}  AuthResponse
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

	domainUser := &domain.User{
		Username: req.Username,
		Password: req.Password,
	}

	registerResponse, err := h.userService.RegisterUser(r.Context(), domainUser)
	if err != nil {
		http.Error(w, "failed to register user", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, registerResponse.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&AuthResponse{
		UserId:   registerResponse.UserID,
		JwtToken: registerResponse.AccessToken,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns access token. Refresh token is set as httpOnly cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Login request"
// @Success      200   {object}  AuthResponse
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

	loginResponse, err := h.userService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "failed to login user", http.StatusInternalServerError)
		return
	}

	h.setRefreshTokenCookie(w, loginResponse.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&AuthResponse{
		UserId:   loginResponse.UserID,
		JwtToken: loginResponse.AccessToken,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Returns new access token. Reads refresh token from httpOnly cookie.
// @Tags         auth
// @Produce      json
// @Success      200   {object}  AuthResponse
// @Failure      401   {string}  string  "missing refresh token"
// @Failure      500   {string}  string  "failed to refresh token"
// @Router       /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	refreshResponse, err := h.userService.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "failed to refresh token", http.StatusUnauthorized)
		return
	}

	h.setRefreshTokenCookie(w, refreshResponse.RefreshToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&AuthResponse{
		UserId:   refreshResponse.UserID,
		JwtToken: refreshResponse.AccessToken,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Revokes the refresh token from httpOnly cookie and clears it
// @Tags         auth
// @Success      204
// @Failure      401   {string}  string  "missing refresh token"
// @Failure      500   {string}  string  "failed to revoke token"
// @Router       /api/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "missing refresh token", http.StatusUnauthorized)
		return
	}

	err = h.userService.RevokeToken(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "failed to revoke token", http.StatusInternalServerError)
		return
	}

	h.clearRefreshTokenCookie(w)
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
