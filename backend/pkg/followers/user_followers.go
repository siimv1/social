package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"strconv"

	"github.com/gorilla/mux"
)

func GetUserFollowersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Retrieve followers from the database
	rows, err := db.DB.Query(`
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM followers f
        JOIN users u ON u.id = f.follower_id
        WHERE f.followed_id = ?`, userID)
	if err != nil {
		log.Printf("Error fetching followers: %v", err)
		http.Error(w, "Failed to fetch followers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []Follower
	for rows.Next() {
		var follower Follower
		if err := rows.Scan(&follower.ID, &follower.FirstName, &follower.LastName, &follower.Nickname); err != nil {
			log.Printf("Error scanning follower: %v", err)
			http.Error(w, "Failed to read followers", http.StatusInternalServerError)
			return
		}
		followers = append(followers, follower)
	}

	// Return followers in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"followers": followers})
}
