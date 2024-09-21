package notifications

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

// Notification represents a notification record.
type Notification struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Read      bool   `json:"read"`
}

// HandleGetUnreadNotifications handles fetching unread notifications for a user.
func HandleGetUnreadNotifications(w http.ResponseWriter, r *http.Request) {
	// Fetch user_id from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Open database connection
	db, err := sql.Open("sqlite3", "./social_network.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Fetch unread notifications
	unreadNotifications, err := GetUnreadNotifications(db, userID)
	if err != nil {
		http.Error(w, "Error fetching notifications", http.StatusInternalServerError)
		return
	}

	// Return notifications as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unreadNotifications)
}

// HandleMarkNotificationAsRead handles marking a notification as read.
func HandleMarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	// Parse notification ID from URL parameters
	notificationIDStr := r.URL.Query().Get("id")
	notificationID, err := strconv.Atoi(notificationIDStr)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	// Open database connection
	db, err := sql.Open("sqlite3", "./social_network.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Mark notification as read
	err = MarkNotificationAsRead(db, notificationID)
	if err != nil {
		http.Error(w, "Error marking notification as read", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
}

// GetUnreadNotifications retrieves unread notifications for a user.
func GetUnreadNotifications(db *sql.DB, userID int) ([]Notification, error) {
	query := `SELECT id, user_id, message, type, created_at, read FROM notifications WHERE user_id = ? AND read = 0`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.UserID, &notification.Message, &notification.Type, &notification.CreatedAt, &notification.Read); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// MarkNotificationAsRead marks a specific notification as read.
func MarkNotificationAsRead(db *sql.DB, notificationID int) error {
	query := `UPDATE notifications SET read = 1 WHERE id = ?`
	_, err := db.Exec(query, notificationID)
	return err
}
