package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"futuremarket/service"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles registration, login, and logout
type AuthHandler struct {
	Service          service.UserService
	BlacklistService service.BlacklistService // REQUIRED FOR LOGOUT
}

// -----------------------------------------------
// POST /api/v1/register
// -----------------------------------------------
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Parse JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// -----------------------------
	// VALIDATION
	// -----------------------------
	if err := service.ValidateName(req.Name); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.ValidateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Role != "user" && req.Role != "admin" {
		http.Error(w, "Role can either be user or admin", http.StatusBadRequest)
		return
	}
	// -----------------------------
	// CREATE USER (service handles hashing + checking duplicates)
	// -----------------------------
	_, err := h.Service.RegisterUser(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// -----------------------------
	// SUCCESS RESPONSE
	// -----------------------------
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "registration successful",
	})
}

// -----------------------------------------------
// POST /api/v1/login
// -----------------------------------------------
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	// Lookup user
	user, err := h.Service.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// JWT claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Successful login → return JWT
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": signedToken,
	})
}

// -----------------------------------------------
// POST /api/v1/logout
// -----------------------------------------------
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Blacklist the token
	err := h.BlacklistService.BlacklistToken(token)
	if err != nil {
		http.Error(w, "failed to blacklist token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successful",
	})
}
