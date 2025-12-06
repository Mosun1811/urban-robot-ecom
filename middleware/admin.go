package middleware

import (
	"fmt"
	"net/http"
)

func AdminMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        roleValue := r.Context().Value(ContextRole) // <-- FIX: use ctxKey, not string
        role, ok := roleValue.(string)

        fmt.Println("ADMIN MIDDLEWARE → ROLE:", role)

        if !ok || role != "admin" {
            http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
    })
}
