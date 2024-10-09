package followers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
)

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

type FollowResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the token
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
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

	// Check if the user is already following or has a pending request
	var existingStatus string
	err := db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followReq.FollowedID).Scan(&existingStatus)
	if err == sql.ErrNoRows {
		// No existing relationship, proceed to follow
	} else if err != nil {
		// Some other error occurred
		log.Printf("Error checking existing follow relationship: %v", err)
		http.Error(w, "Failed to check follow relationship", http.StatusInternalServerError)
		return
	} else {
		// Follow relationship already exists
		log.Printf("User %d already has a follow relationship with user %d", userID, followReq.FollowedID)
		http.Error(w, "Follow request already exists", http.StatusConflict)
		return
	}

	// Check if the followed user is public or private
	var isPublic int
	err = db.DB.QueryRow("SELECT is_public FROM users WHERE id = ?", followReq.FollowedID).Scan(&isPublic)
	if err != nil {
		log.Printf("Error checking user visibility: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	log.Printf("Followed user (ID: %d) is_public value: %d", followReq.FollowedID, isPublic)

	// Set the follow status to 'pending' if private or 'accepted' if public
	var status string
	if isPublic == 1 {
		status = "accepted"
	} else {
		status = "pending"
	}

	// Insert the follow relationship with the appropriate status
	_, err = db.DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", userID, followReq.FollowedID, status)
	if err != nil {
		log.Printf("Error following user %d: %v", followReq.FollowedID, err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	// Success response
	var message string
	if status == "accepted" {
		message = "You are now following the user."
	} else {
		message = "Follow request sent and is pending approval."
	}

	resp := FollowResponse{
		Message: message,
		Status:  status,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the token
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
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

	_, err := db.DB.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", userID, unfollowReq.FollowedID)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	// Success response
	resp := map[string]string{"status": "OK", "message": "Unfollowed successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
