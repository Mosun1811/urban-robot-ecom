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
	"net/http"

	"gorm.io/gorm"
)

// AuthHandler handles registration and login (Epic 1: IAM).
type AuthHandler struct {
	DB *gorm.DB
}

// POST /api/v1/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement user registration"))
}

// POST /api/v1/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("TODO: implement user login"))
}
