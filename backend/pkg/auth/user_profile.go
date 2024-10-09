package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"strconv"
	"strings"
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
	// Extract the user ID from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	userIDStr := pathParts[len(pathParts)-1]
	profileID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get the current logged-in user ID from the token
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	currentUserID, err := ValidateToken(tokenString)
	if err != nil {
		currentUserID = 0 // Not logged in
	}

	// Fetch user profile data
	var user User
	err = db.DB.QueryRow("SELECT id, first_name, last_name, nickname, email, date_of_birth, about_me, avatar, is_public FROM users WHERE id = ?", profileID).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Email, &user.DateOfBirth, &user.AboutMe, &user.Avatar, &user.IsPublic)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Determine follow status
	var followStatus string = "not-following"
	if currentUserID != 0 && currentUserID != profileID {
		err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", currentUserID, profileID).Scan(&followStatus)
		if err != nil {
			followStatus = "not-following"
		}
	}

	// Prepare the response data
	response := map[string]interface{}{
		"id":            user.ID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"nickname":      user.Nickname,
		"email":         user.Email,
		"date_of_birth": user.DateOfBirth,
		"about_me":      user.AboutMe,
		"avatar":        user.Avatar,
		"is_public":     user.IsPublic, // It's already a boolean, no need to compare with 1 or 0
		"follow_status": followStatus,
		"is_private":    !user.IsPublic, // Simply negate the boolean to check if the profile is private
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func IsFollower(followerID, followedID int) (bool, error) {
	var status string
	err := db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, followedID).Scan(&status)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return status == "accepted", nil
}

func IsFollowing(userID, followedID int) (bool, error) {
	return IsFollower(userID, followedID)
}
