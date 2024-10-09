package following

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

func GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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

	rows, err := db.DB.Query(`
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM followers f
        JOIN users u ON u.id = f.followed_id
        WHERE f.follower_id = ? AND f.status = 'accepted'`, userID)
	if err != nil {
		log.Printf("Error fetching following users: %v", err)
		http.Error(w, "Failed to fetch following users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var following []Follower
	for rows.Next() {
		var follower Follower
		if err := rows.Scan(&follower.ID, &follower.FirstName, &follower.LastName, &follower.Nickname); err != nil {
			log.Printf("Error scanning following user: %v", err)
			http.Error(w, "Failed to read following users", http.StatusInternalServerError)
			return
		}
		following = append(following, follower)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"following": following})
}
