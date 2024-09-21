package followers

import (
	"encoding/json"
	"net/http"
	"social-network/backend/pkg/db"
)

// getFollowingHandler: Järgjate saamine
func GetFollowingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	rows, err := db.DB.Query("SELECT followed_id FROM followers WHERE follower_id = ?", followerID)
	if err != nil {
		http.Error(w, "Failed to retrieve following users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followingUsers []User
	for rows.Next() {
		var followedID int
		if err := rows.Scan(&followedID); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}

		var user User
		err := db.DB.QueryRow("SELECT id, first_name, last_name FROM users WHERE id = ?", followedID).Scan(&user.ID, &user.FirstName, &user.LastName)
		if err != nil {
			continue // Ignore errors for individual users
		}

		user.IsFollowing = true
		followingUsers = append(followingUsers, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": followingUsers,
	})
}
