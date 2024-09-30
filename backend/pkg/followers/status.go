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
	// URL-st võtame ID
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

	// Kontrollime jälgimise staatust
	var status string
	err = db.DB.QueryRow("SELECT status FROM followers WHERE follower_id = ? AND followed_id = ?", userID, followedID).Scan(&status)
	if err != nil {
		log.Printf("Error checking follow status: %v", err)
		http.Error(w, "Failed to check follow status", http.StatusInternalServerError)
		return
	}

	// Tagastame jälgimise staatuse kliendile
	if status == "following" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "following"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "not-following"}`))
	}
}
