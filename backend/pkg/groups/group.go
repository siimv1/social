package groups

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
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

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	var creatorID int
	err = db.DB.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)
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
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	query := `INSERT INTO join_requests (group_id, user_id, status) VALUES (?, ?, 'pending')`
	_, err = db.DB.Exec(query, groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func AcceptJoinRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	var groupID, userID int
	err := db.DB.QueryRow("SELECT group_id, user_id FROM join_requests WHERE id = ?", requestID).Scan(&groupID, &userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("UPDATE join_requests SET status = 'accepted' WHERE id = ?", requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RejectJoinRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	query := `UPDATE join_requests SET status = 'rejected' WHERE id = ?`
	_, err := db.DB.Exec(query, requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// func GetJoinStatus(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	groupID := vars["id"]

// 	tokenStr := r.Header.Get("Authorization")
// 	if tokenStr == "" {
// 		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
// 		return
// 	}

// 	userID, err := auth.ValidateToken(tokenStr)
// 	if err != nil {
// 		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	var isPending, isMember bool

// 	err = db.DB.QueryRow("SELECT COUNT(*) > 0 FROM join_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'", groupID, userID).Scan(&isPending)
// 	if err != nil {
// 		log.Printf("Error checking pending join requests: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	err = db.DB.QueryRow("SELECT COUNT(*) > 0 FROM group_members WHERE group_id = ? AND user_id = ?", groupID, userID).Scan(&isMember)
// 	if err != nil {
// 		log.Printf("Error checking if user is member: %v", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]bool{
// 		"isPending": isPending,
// 		"isMember":  isMember,
// 	})
// }

func JoinGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]
	tokenStr := r.Header.Get("Authorization")

	userID, err := auth.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var existingRequestCount int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM join_requests WHERE group_id = ? AND user_id = ? AND status = 'pending'", groupID, userID).Scan(&existingRequestCount)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existingRequestCount > 0 {
		http.Error(w, "Join request already sent", http.StatusBadRequest)
		return
	}

	// Salvestame uue liitumisavalduse
	_, err = db.DB.Exec("INSERT INTO join_requests (group_id, user_id, status) VALUES (?, ?, 'pending')", groupID, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}