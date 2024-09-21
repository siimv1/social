package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
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

type FollowResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// FollowHandler: Järgimise käsitlemine
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

	log.Printf("Following user ID: %d by follower ID: %d", followReq.FollowedID, followerID)

	var existingStatus string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, followReq.FollowedID).Scan(&existingStatus)
	if err == nil {
		http.Error(w, "Already following", http.StatusConflict)
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

	_, err = db.DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", followerID, followReq.FollowedID, status)
	if err != nil {
		log.Printf("Error following user: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	resp := FollowResponse{
		Message: "Follow request sent",
		Status:  status,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UnfollowHandler: Üks järgimine lõpetamine
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

	var unfollowReq UnfollowRequest
	if err := json.NewDecoder(r.Body).Decode(&unfollowReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, unfollowReq.FollowedID)
	if err != nil {
		log.Printf("Error unfollowing user: %v", err)
		http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "Unfollowed successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
