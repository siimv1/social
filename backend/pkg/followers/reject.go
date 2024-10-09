package followers

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
)

func RejectFollowRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Get the recipient ID from the context
	recipientID, ok := r.Context().Value(auth.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body to get the follower ID
	var req struct {
		FollowerID int `json:"follower_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the follow request status to 'rejected'
	result, err := db.DB.Exec("UPDATE followers SET status = 'rejected' WHERE follower_id = ? AND followed_id = ? AND status = 'pending'", req.FollowerID, recipientID)
	if err != nil {
		log.Printf("Error rejecting follow request: %v", err)
		http.Error(w, "Failed to reject follow request", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No pending follow request found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Follow request rejected"})
}
