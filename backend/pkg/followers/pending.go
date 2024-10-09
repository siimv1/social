package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
)

func GetPendingFollowRequestsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the recipient ID from the context
	recipientID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Fetching pending follow requests for user ID: %d", recipientID)

	// Fetch pending follow requests
	rows, err := db.DB.Query(`
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM followers f
        JOIN users u ON u.id = f.follower_id
        WHERE f.followed_id = ? AND f.status = 'pending'`, recipientID)
	if err != nil {
		log.Printf("Error fetching pending follow requests: %v", err)
		http.Error(w, "Failed to fetch pending follow requests", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []Follower
	for rows.Next() {
		var follower Follower
		if err := rows.Scan(&follower.ID, &follower.FirstName, &follower.LastName, &follower.Nickname); err != nil {
			log.Printf("Error scanning follower: %v", err)
			http.Error(w, "Failed to read pending follow requests", http.StatusInternalServerError)
			return
		}
		requests = append(requests, follower)
	}

	// Return the pending follow requests
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"requests": requests})
}
