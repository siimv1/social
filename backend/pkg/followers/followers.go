package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
)

type FollowRequest struct {
	FollowedID int `json:"followed_id"`
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userEmail := r.Header.Get("User-Email")
	if userEmail == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var followerID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", userEmail).Scan(&followerID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	var followReq FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var followedPrivacy string
	err = db.DB.QueryRow("SELECT privacy FROM users WHERE id = ?", followReq.FollowedID).Scan(&followedPrivacy)
	if err != nil {
		http.Error(w, "User to follow not found", http.StatusNotFound)
		return
	}
	status := "pending"
	if followedPrivacy == "public" {
		status = "accepted"
	}
	query := `INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)`
	_, err = db.DB.Exec(query, followerID, followReq.FollowedID, status)
	if err != nil {
		log.Printf("Error following user: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request sent"})
}

// UnfollowHandler allows a user to unfollow another user
func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userEmail := r.Header.Get("User-Email")
	if userEmail == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var followerID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = ?", userEmail).Scan(&followerID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	var followReq FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	// Delete the follow relationship
	query := `DELETE FROM followers WHERE follower_id = ? AND followed_id = ?`
	_, err = db.DB.Exec(query, followerID, followReq.FollowedID)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Unfollowed successfully"})
}
