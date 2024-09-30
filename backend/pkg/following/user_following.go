package following

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"strconv"

	"github.com/gorilla/mux"
)

type Following struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
}

func GetUserFollowingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Retrieve following users
	rows, err := db.DB.Query(`
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM followers f
        JOIN users u ON u.id = f.followed_id
        WHERE f.follower_id = ?`, userID)
	if err != nil {
		log.Printf("Error fetching following users: %v", err)
		http.Error(w, "Failed to fetch following users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var following []Following
	for rows.Next() {
		var user Following
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname); err != nil {
			log.Printf("Error scanning following user: %v", err)
			http.Error(w, "Failed to read following users", http.StatusInternalServerError)
			return
		}
		following = append(following, user)
	}

	// Return following users in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"following": following})
}
