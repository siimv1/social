package followers

import (
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
	"strconv"
	"strings"
)

func CheckFollowStatusHandler(w http.ResponseWriter, r *http.Request) {
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

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userID, err := auth.ValidateToken(tokenString)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var status string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followedID).Scan(&status)
	if err != nil {
		log.Printf("Error checking follow status: %v", err)
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}

	if status == "following" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "following"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "not-following"}`))
	}
}
// For chat
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
