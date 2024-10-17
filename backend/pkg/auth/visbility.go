package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
)

// UpdateProfileVisibilityHandler handles requests to update the user's profile visibility.
func UpdateProfileVisibilityHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the session
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Println("Failed to get session:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userIDInterface, ok := session.Values["user_id"]
	if !ok {
		log.Println("User ID not found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userID int
	switch v := userIDInterface.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	default:
		log.Printf("Invalid user ID type in session: %T", v)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the request body to get the new visibility setting
	var visibilityReq struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := json.NewDecoder(r.Body).Decode(&visibilityReq); err != nil {
		log.Printf("Invalid request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert bool to int for database storage
	isPublicInt := 0
	if visibilityReq.IsPublic {
		isPublicInt = 1
	}

	// Update the user's profile visibility in the database
	result, err := db.DB.Exec("UPDATE users SET is_public = ? WHERE id = ?", isPublicInt, userID)
	if err != nil {
		log.Printf("Failed to update profile visibility for user %d: %v", userID, err)
		http.Error(w, "Failed to update profile visibility", http.StatusInternalServerError)
		return
	}

	// Check if the update affected any rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
	} else {
		log.Printf("Rows affected: %d", rowsAffected)
	}

	// Return a success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"isPublic": visibilityReq.IsPublic,
	})
}
