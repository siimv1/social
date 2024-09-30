package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
)

// Follower defines the structure for a follower
type Follower struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
}

func GetFollowersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve the user ID from the context
	userIDValue := r.Context().Value(auth.UserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized: no user ID in context", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Unauthorized: invalid user ID type", http.StatusUnauthorized)
		return
	}

	// Log after confirming userID
	log.Println("Fetching followers for userID:", userID)

	// Retrieve followers from the database using the userID
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

	// Return the followers in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"followers": followers})
}
