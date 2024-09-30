package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
	"strings"
)

// Järgijate ja kasutaja defineerimine
type User struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	IsFollowing bool   `json:"is_following"`
}

type FollowRequest struct {
	FollowedID int `json:"followed_id"`
}

type UnfollowRequest struct {
	FollowedID int `json:"followed_id"`
}

// FollowResponse on vastuse struktuur
type FollowResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// followers/follow.go
func FollowHandler(w http.ResponseWriter, r *http.Request) {
	// Get the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Trim the 'Bearer ' prefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Validate the token and get the user ID
	userID, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Read the request body to get followed_id
	var followReq FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followReq); err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check if the user is already following
	var existingStatus string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followReq.FollowedID).Scan(&existingStatus)
	if err == nil {
		log.Printf("User %d is already following user %d", userID, followReq.FollowedID)
		http.Error(w, "Already following", http.StatusConflict)
		return
	}

	// Insert the follow relationship
	_, err = db.DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", userID, followReq.FollowedID, "accepted")
	if err != nil {
		log.Printf("Error following user %d: %v", followReq.FollowedID, err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	// Success
	resp := FollowResponse{
		Message: "Follow request successful",
		Status:  "accepted",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Get the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Trim the 'Bearer ' prefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	// Validate the token and get the user ID
	userID, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Read the request body to get followed_id
	var unfollowReq UnfollowRequest
	if err := json.NewDecoder(r.Body).Decode(&unfollowReq); err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", userID, unfollowReq.FollowedID)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}
	// Korrektne vastus, mida frontendi poolel oodatakse
	resp := map[string]string{"status": "OK", "message": "Unfollowed successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
