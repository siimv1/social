package groups

import (
	"encoding/json"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"

	"github.com/gorilla/mux"
)

func GetEventInvites(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.DB.Query(`
        SELECT eventinvites.id, eventinvites.event_id, events.title
        FROM eventinvites
        JOIN events ON eventinvites.event_id = events.id
        WHERE eventinvites.user_id = ? AND eventinvites.status = 'pending'`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type EventInvite struct {
		ID         int    `json:"id"`
		EventID    int    `json:"event_id"`
		EventTitle string `json:"event_title"`
	}

	var eventInvites []EventInvite

	for rows.Next() {
		var invite EventInvite
		if err := rows.Scan(&invite.ID, &invite.EventID, &invite.EventTitle); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		eventInvites = append(eventInvites, invite)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"eventinvites": eventInvites,
	})
}

func AcceptEventInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteID := vars["inviteId"]

	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := db.DB.Exec("UPDATE eventinvites SET status = 'accepted' WHERE id = ? AND user_id = ?", inviteID, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var eventID int
	err = db.DB.QueryRow("SELECT event_id FROM eventinvites WHERE id = ?", inviteID).Scan(&eventID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO event_participants (event_id, user_id) VALUES (?, ?)", eventID, userID)
	if err != nil {
		http.Error(w, "Failed to add user to event participants", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeclineEventInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteID := vars["inviteId"]

	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := db.DB.Exec("UPDATE eventinvites SET status = 'declined' WHERE id = ? AND user_id = ?", inviteID, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func InviteUserToEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	session, _ := auth.Store.Get(r, "session-name")
	currentUserID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var creatorID int
	err := db.DB.QueryRow("SELECT creator_id FROM events WHERE id = ?", eventID).Scan(&creatorID)
	if err != nil || creatorID != currentUserID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var requestData struct {
		UserID int `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("INSERT INTO eventinvites (event_id, user_id) VALUES (?, ?)", eventID, requestData.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var eventData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Date        string `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&eventData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("INSERT INTO events (title, description, date, creator_id) VALUES (?, ?, ?, ?)",
		eventData.Title, eventData.Description, eventData.Date, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve event ID", http.StatusInternalServerError)
		return
	}

	response := map[string]int64{
		"event_id": eventID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
