package auth

import (
	"encoding/json"
	"net/http"
	"social-network/backend/pkg/db"
	"time"

	"github.com/golang-jwt/jwt" 
	"golang.org/x/crypto/bcrypt"
)

// Secret key used to sign tokens (securely store this in production)
var jwtKey = []byte("secret_key")

// Claims represents the JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	var userID int
	var hashedPassword string
	query := `SELECT id, password FROM users WHERE email = ?`
	err := db.DB.QueryRow(query, loginReq.Email).Scan(&userID, &hashedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Create the JWT claims, which includes the user ID and expiry time
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID, // userID obtained after validating user credentials
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	// Return the token in the response
// Return the token and user ID in the response
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]interface{}{
    "token": tokenString,
    "user_id": userID, // Tagastame ka kasutaja ID
}) 
}


// LogoutHandler logs the user out by clearing the JWT token
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour), // Expires immediately
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}