package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"strings"

	"github.com/golang-jwt/jwt"
)

func GetUserProfile(userID int) (*User, error) {
	var user User
	query := `SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, is_public FROM users WHERE id = ?`
	err := db.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Avatar,
		&user.Nickname,
		&user.AboutMe,
		&user.IsPublic, // Include this field
	)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		return nil, err
	}
	return &user, nil
}

// ProfileHandler handles requests to get the logged-in user's profile
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the JWT token from the Authorization header
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Remove "Bearer " prefix from the token
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

	// Parse the JWT token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Check if the token is valid
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch user profile based on the userID from the claims
	user, err := GetUserProfile(claims.UserID)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Exclude the password from the response
	user.Password = ""

	// Return the user profile as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
