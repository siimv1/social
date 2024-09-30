package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type UserProfile struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
	AboutMe     string `json:"about_me"`
	Avatar      string `json:"avatar"`
	IsFollowing bool   `json:"is_following"`
}

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Kontrolli, kas päringul on kehtiv token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	requesterID, err := ValidateToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Toome kasutaja andmed
	var user UserProfile
	err = db.DB.QueryRow(`
        SELECT id, first_name, last_name, nickname, email, date_of_birth, about_me, avatar
        FROM users
        WHERE id = ?
    `, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Email, &user.DateOfBirth, &user.AboutMe, &user.Avatar)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Kontrollime, kas requester juba jälgib
	var isFollowing bool
	err = db.DB.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM followers 
            WHERE follower_id = ? AND followed_id = ? AND status = 'accepted'
        )
    `, requesterID, userID).Scan(&isFollowing)
	if err != nil {
		log.Printf("Error checking follow status: %v", err)
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}
	user.IsFollowing = isFollowing

	// Tagastame kasutaja andmed JSON formaadis
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
