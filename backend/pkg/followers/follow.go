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
	Message     string `json:"message"`
	Status      string `json:"status"`
	IsFollowing bool   `json:"is_following"`
}

// FollowHandler handles following a user
func FollowHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("No user_id found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var followReq FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&followReq); err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Kontrolli, kas juba on olemasolev suhe (accepted või pending)
	var existingStatus string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followReq.FollowedID).Scan(&existingStatus)
	if err == sql.ErrNoRows {
		// Suhet pole, jätka jälgimist
	} else if err != nil {
		log.Printf("Error checking existing follow relationship: %v", err)
		http.Error(w, "Failed to check follow relationship", http.StatusInternalServerError)
		return
	} else {
		// Kui juba on suhe (accepted või pending), tagasta veateade
		if existingStatus == "accepted" || existingStatus == "pending" {
			log.Printf("User %d already has a follow relationship with user %d", userID, followReq.FollowedID)
			http.Error(w, "Follow request already exists", http.StatusConflict)
			return
		}
	}

	// Kontrolli, kas kasutaja on avalik või privaatne
	var isPublic int
	err = db.DB.QueryRow("SELECT is_public FROM users WHERE id = ?", followReq.FollowedID).Scan(&isPublic)
	if err != nil {
		log.Printf("Error checking user visibility: %v", err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	status := "pending"
	if isPublic == 1 {
		status = "accepted"
	}

	// Lisa uus jälgimissuhe andmebaasi
	_, err = db.DB.Exec("INSERT INTO followers (follower_id, followed_id, status) VALUES (?, ?, ?)", userID, followReq.FollowedID, status)
	if err != nil {
		log.Printf("Error following user %d: %v", followReq.FollowedID, err)
		http.Error(w, "Failed to follow user", http.StatusInternalServerError)
		return
	}

	var followStatus string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followReq.FollowedID).Scan(&followStatus)
	if err == sql.ErrNoRows {
		followStatus = "not-following"
	} else if err != nil {
		log.Printf("Error checking follow status: %v", err)
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}

	isFollowing := (followStatus == "accepted")

	resp := FollowResponse{
		Message:     "Follow status checked",
		Status:      followStatus,
		IsFollowing: isFollowing,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UnfollowHandler handles unfollowing a user
func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	session, err := auth.Store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("No user_id found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	resp := map[string]string{"status": "OK", "message": "Unfollowed successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
