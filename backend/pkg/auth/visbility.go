package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
)

func UpdateProfileVisibilityHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var visibilityReq struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&visibilityReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert bool to int for SQLite (1 for true, 0 for false)
	isPublicInt := 0
	if visibilityReq.IsPublic {
		isPublicInt = 1
	}

	_, err := db.DB.Exec("UPDATE users SET is_public = ? WHERE id = ?", isPublicInt, userID)
	if err != nil {
		log.Printf("Failed to update profile visibility for user %d: %v", userID, err)
		http.Error(w, "Failed to update profile visibility", http.StatusInternalServerError)
		return
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"isPublic": visibilityReq.IsPublic,
	})
}
