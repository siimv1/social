package auth

import (
	"context"
	"log"
	"net/http"
)

// ContextKey is a custom type to prevent context key collisions
type ContextKey string

// UserIDKey is the key used to store the user ID in the context
const UserIDKey ContextKey = "userID"

// AuthMiddleware checks if the user is authenticated via session cookies
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, "session-name")
		if err != nil {
			log.Println("Failed to get session:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userIDInterface, ok := session.Values["user_id"]
		if !ok {
			log.Println("User ID not found in session")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var userID int
		switch v := userIDInterface.(type) {
		case int:
			userID = v
		case int64:
			userID = int(v)
		default:
			log.Printf("Invalid user ID type in session: %T", v)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
