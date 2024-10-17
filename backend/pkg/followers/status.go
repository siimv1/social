package followers

import (
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
	"strconv"
	"strings"
)

// CheckFollowStatusHandler checks if the current user is following another user
func CheckFollowStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the followed user ID from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	followedIDStr := pathParts[len(pathParts)-1]
	followedID, err := strconv.Atoi(followedIDStr)
	if err != nil {
		http.Error(w, "Invalid followed ID", http.StatusBadRequest)
		return
	}

	// Get the current logged-in user ID from the session (set by the AuthMiddleware)
	userID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check follow status in the database
	var status string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followedID).Scan(&status)
	if err != nil {
		log.Printf("Error checking follow status: %v", err)
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}

	// Respond with the follow status
	if status == "accepted" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "following"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "not-following"}`))
	}
}

// CheckMutualFollowStatus checks if two users follow each other
func CheckMutualFollowStatus(followerID, followedID int) (bool, error) {
	var status string

	log.Printf("Checking if user %d is following user %d", followerID, followedID)
	// Check if the follower is following the followed user
	err := db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, followedID).Scan(&status)
	if err != nil {
		log.Printf("User %d is not following user %d: %v", followerID, followedID, err)
	} else if status == "accepted" {
		log.Printf("User %d is following user %d with status: %s", followerID, followedID, status)
		return true, nil
	}

	log.Printf("Checking if user %d is following user %d", followedID, followerID)
	// Check if the followed user is following back
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", followedID, followerID).Scan(&status)
	if err != nil {
		log.Printf("User %d is not following user %d: %v", followedID, followerID, err)
		return false, nil
	}

	if status == "accepted" {
		log.Printf("User %d is following user %d with status: %s", followedID, followerID, status)
		return true, nil
	}

	log.Printf("Users %d and %d are not following each other", followerID, followedID)
	return false, nil
}
