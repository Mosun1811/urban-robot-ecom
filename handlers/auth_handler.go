// This file is assigned to Munson, for registering and logging in users.
// These are important without it users wont be able to access the protected parts of our API.
//
// PURPOSE:
// - HTTP handlers for authentication and identity (register + login).
// - Entry point for creating users and issuing JWT tokens.
//
// EPICS & USER STORIES: (Here i've mapped teh relevant epics and user stories so you
// can reference either JIRA or the Project reqirements from Oreva)
// - Epic 1: Identity & Access Management (IAM)
//   - User Story 1.1: User Registration  (POST /api/v1/register)
//   - User Story 1.2: Authentication     (POST /api/v1/login)
//
// ENDPOINTS (These are all the endpoints we need to include here):
// - POST /api/v1/register
//   - Accepts: { name, email, password }
//   - Validates input, hashes password (bcrypt), stores user in DB.
//   - Returns 404.
//
// - POST /api/v1/login
//   - Accepts: { email, password }
//   - Verifies password hash.
//   - Returns JWT with user_id + role (admin/customer), 24h expiry.

//( here we are just making sure that when a user logs in we check if the email exists, compare the password
// the hased one, and if everything matches, we create a JWT token.
//That token includes user's id and role (customer/admin) and is what we'll use to protect all cart/order/review routes etc)

// Hope this makes sense, let me know!!

// What I have done below is just to build so that everything compiles and you'll be able to clone have working code
// Only thing you'd need to do is to write the logic

package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"futuremarket/models"
	"futuremarket/service"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles registration and login
type AuthHandler struct {
	Service service.UserService
}

// -----------------------------------------------
// POST /api/v1/register
// -----------------------------------------------
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse input JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email and password required", http.StatusBadRequest)
		return
	}

	// Check if email already exists
	_, err := h.Service.GetUserByEmail(req.Email)
	if err == nil {
		http.Error(w, "user with email already exists", http.StatusBadRequest)
		return
	}
	
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	
	// Create user
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "customer",
	}

	err = h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	// Response
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

	// Parse JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password required", http.StatusBadRequest)
		return
	}

	// Find user
	user, err := h.Service.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare password
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

	// Send back the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": signedToken,
	})
}

// -----------------------------------------------
// POST /api/v1/logout
// -----------------------------------------------
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successful — please delete your token on the client side",
	})
}
