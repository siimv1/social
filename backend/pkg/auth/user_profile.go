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

// UserProfile contains user profile details
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

// UserProfileHandler handles requests for viewing user profiles
// Fetch user profile data, including follow status
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get logged-in user's ID from the session
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	loggedInUserID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("No user_id found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract the profile ID from the URL path
	profileID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Fetch user profile
	var user User
	err = db.DB.QueryRow("SELECT id, first_name, last_name, nickname, email, date_of_birth, about_me, avatar, is_public FROM users WHERE id = ?", profileID).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname, &user.Email, &user.DateOfBirth, &user.AboutMe, &user.Avatar, &user.IsPublic)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the logged-in user is following the profile
	var followStatus string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", loggedInUserID, profileID).Scan(&followStatus)
	if err == sql.ErrNoRows {
		followStatus = "not-following"
	} else if err != nil {
		log.Printf("Error checking follow status: %v", err)
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}

	// Include follow_status in the response
	response := map[string]interface{}{
		"id":            user.ID,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"nickname":      user.Nickname,
		"email":         user.Email,
		"date_of_birth": user.DateOfBirth,
		"about_me":      user.AboutMe,
		"avatar":        user.Avatar,
		"is_public":     user.IsPublic,
		"follow_status": followStatus, // Tagasta follow status
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// IsFollower checks if the user is following the given profile
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

// IsFollowing is an alias for IsFollower, kept for semantic clarity
func IsFollowing(userID, followedID int) (bool, error) {
	return IsFollower(userID, followedID)
}

// SomeProtectedHandler is an example of a protected endpoint that returns user-specific data
func SomeProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the context
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Use the userID to fetch user-specific data from the database
	userData, err := fetchUserData(userID)
	if err != nil {
		log.Printf("Failed to fetch user data for user ID %d: %v", userID, err)
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	// Return the user data as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}

// fetchUserData retrieves user-specific data from the database
func fetchUserData(userID int) (UserData, error) {
	var userData UserData
	err := db.DB.QueryRow(`
        SELECT id, first_name, last_name, email, date_of_birth, about_me
        FROM users
        WHERE id = ?
    `, userID).Scan(
		&userData.ID,
		&userData.FirstName,
		&userData.LastName,
		&userData.Email,
		&userData.DateOfBirth,
		&userData.AboutMe,
	)
	if err != nil {
		return UserData{}, err
	}
	return userData, nil
}

// UserData represents the user's data structure
type UserData struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateOfBirth string `json:"date_of_birth"`
	AboutMe     string `json:"about_me"`
}
