package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
)

// GetUserProfile fetches the user profile by user ID
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
		&user.IsPublic,
	)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		return nil, err
	}
	return &user, nil
}

// ProfileHandler handles requests to get the logged-in user's profile
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the context (set by the AuthMiddleware)
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch the user profile using the user ID from the session
	user, err := GetUserProfile(userID)
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
