package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// This interface allows your service to be injected cleanly.
type BlacklistService interface {
	IsTokenBlacklisted(token string) (bool, error)
}

type AuthMiddlewareConfig struct {
	BlacklistService BlacklistService
}

type ctxKey string

const (
	ContextUserID ctxKey = "user_id"
	ContextRole   ctxKey = "role"
)

func (cfg AuthMiddlewareConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Read Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. 🔥 Blacklist check BEFORE parsing JWT
		if cfg.BlacklistService != nil {
			isBlacklisted, err := cfg.BlacklistService.IsTokenBlacklisted(tokenStr)
			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}

			if isBlacklisted {
				http.Error(w, "token is blacklisted", http.StatusUnauthorized)
				return
			}
		}

		// 3. Parse and validate JWT
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		fmt.Println("ROLE CLAIM FROM JWT:", claims["role"])       // ← ADD THIS
		fmt.Println("USER ID CLAIM FROM JWT:", claims["user_id"]) // optional

		// user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "invalid user_id claim", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		// role
		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "invalid role claim", http.StatusUnauthorized)
			return
		}

		fmt.Println("AUTH MIDDLEWARE → Role claim:", role)

		// 5. Store values in context
		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextRole, role)

		// 6. Continue the request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
