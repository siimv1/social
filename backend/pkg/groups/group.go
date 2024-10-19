package groups

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Group struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int       `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
}
type Invite struct {
	ID      int    `json:"id"`
	GroupID int    `json:"group_id"`
	Email   string `json:"email"`
	Status  string `json:"status"`
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}
	log.Printf("Creator ID: %d", userID)
	var group Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	group.CreatorID = userID
	query := `INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?)`
	result, err := db.DB.Exec(query, group.Title, group.Description, group.CreatorID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error retrieving last insert ID: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	group.ID = int(id)
	group.CreatedAt = time.Now()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}
func GetGroups(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, description, creator_id, created_at FROM groups")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var groups []Group
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID, &group.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		groups = append(groups, group)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
func GetGroupByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var group Group
	err := db.DB.QueryRow("SELECT id, title, description, creator_id, created_at FROM groups WHERE id = ?", id).Scan(&group.ID, &group.Title, &group.Description, &group.CreatorID, &group.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Group not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var users []struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	for rows.Next() {
		var user struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func InviteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	var creatorID int
	err := db.DB.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Group not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if creatorID != userID {
		http.Error(w, "You are not the creator of this group", http.StatusForbidden)
		return
	}

	var invite Invite
	if err := json.NewDecoder(r.Body).Decode(&invite); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO invites (group_id, email, status) VALUES (?, ?, 'pending')`
	_, err = db.DB.Exec(query, groupID, invite.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func AcceptInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteID := vars["id"]
	query := `UPDATE invites SET status = 'accepted' WHERE id = ?`
	_, err := db.DB.Exec(query, inviteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func RequestJoinGroup(w http.ResponseWriter, r *http.Request) {
	session, _ := auth.Store.Get(r, "session-name")
	userIDValue, exists := session.Values["user_id"]
	if !exists {
		http.Error(w, "Kasutaja ei ole sisse logitud", http.StatusUnauthorized)
		return
	}

	var userID int
	switch v := userIDValue.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	case float64:
		userID = int(v)
	case string:
		var err error
		userID, err = strconv.Atoi(v)
		if err != nil {
			http.Error(w, "Vigane kasutaja ID sessioonis", http.StatusUnauthorized)
			return
		}
	default:
		http.Error(w, "Kasutaja ei ole sisse logitud", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	query := `INSERT INTO join_requests (group_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err := db.DB.Exec(query, groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Liitumisavaldus saadetud gruppi %s", groupID)))
}

func AcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupID := vars["id"]

    var requestData struct {
        UserID int `json:"userId"`
    }
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        log.Printf("Failed to decode request body in AcceptJoinRequest: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("Attempting to accept join request for Group ID: %s, User ID: %d", groupID, requestData.UserID)

    _, err := db.DB.Exec("UPDATE join_requests SET status = 'accepted' WHERE group_id = ? AND user_id = ?", groupID, requestData.UserID)
    if err != nil {
        log.Printf("Failed to execute database update in AcceptJoinRequest: %v", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    log.Printf("Successfully accepted join request for Group ID: %s, User ID: %d", groupID, requestData.UserID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Join request accepted"})
}

func DenyJoinRequest(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupID := vars["id"]

    var requestData struct {
        UserID int `json:"userId"`
    }

    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        log.Printf("Failed to decode request body in DenyJoinRequest: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    log.Printf("Attempting to deny join request for Group ID: %s, User ID: %d", groupID, requestData.UserID)

    _, err := db.DB.Exec("UPDATE join_requests SET status = 'denied' WHERE group_id = ? AND user_id = ?", groupID, requestData.UserID)
    if err != nil {
        log.Printf("Failed to execute database update in DenyJoinRequest: %v", err)
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }

    log.Printf("Successfully denied join request for Group ID: %s, User ID: %d", groupID, requestData.UserID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Join request denied"})
}


func GetJoinRequests(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    groupID := vars["id"]
    
    rows, err := db.DB.Query("SELECT user_id FROM join_requests WHERE group_id = ? AND status = 'pending'", groupID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var requests []struct {
        UserID int `json:"user_id"`
    }
    for rows.Next() {
        var request struct {
            UserID int `json:"user_id"`
        }
        if err := rows.Scan(&request.UserID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        requests = append(requests, request)
    }

    // If no join requests are found, return an empty array
    if len(requests) == 0 {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode([]struct{}{})  // Return an empty array instead of null
        return
    }

    // Send the list of requests as JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(requests)
}


func JoinGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	var creatorID int
	err := db.DB.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if userID == creatorID {
		http.Error(w, "Group creator cannot send join request", http.StatusBadRequest)
		return
	}

	var existingRequestCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM join_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'", groupID, userID).Scan(&existingRequestCount)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingRequestCount > 0 {

		response := map[string]bool{
			"requestPending": true,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		return
	}

	_, err = db.DB.Exec("INSERT INTO join_requests (group_id, user_id, status) VALUES (?, ?, 'pending')", groupID, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{
		"requestPending": true,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func JoinStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	session, _ := auth.Store.Get(r, "session-name")
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	var requestPending bool
	err := db.DB.QueryRow("SELECT COUNT(*) > 0 FROM join_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'", groupID, userID).Scan(&requestPending)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	response := map[string]bool{
		"requestPending": requestPending,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
