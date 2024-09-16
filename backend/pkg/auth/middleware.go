package auth

import (
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// No session cookie found, user is not authenticated
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		sessionToken := cookie.Value
		email, exists := sessionStore[sessionToken]
		if !exists {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// User is authenticated, continue with the request
		r.Header.Set("User-Email", email)
		next(w, r)
	}
}
