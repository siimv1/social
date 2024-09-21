package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"strings"
)

func GetUserProfile(email string) (*User, error) {
	var user User
	query := `SELECT email, first_name, last_name, date_of_birth, avatar, nickname, about_me FROM users WHERE email = ?`
	err := db.DB.QueryRow(query, email).Scan(&user.Email, &user.FirstName, &user.LastName, &user.DateOfBirth, &user.Avatar, &user.Nickname, &user.AboutMe)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Eemaldame "Bearer " tokenist
	sessionToken := strings.TrimPrefix(token, "Bearer ")

	email, exists := sessionStore[sessionToken]
	if !exists {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := GetUserProfile(email)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user) // Saadame kasutaja andmed JSON-formaadis
}
