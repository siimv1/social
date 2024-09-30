package notifications

import (
	"database/sql"
)

// Notification represents a notification entity
type Notification struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
}

// NotificationServe defines methods to interact with notifications
type NotificationServe interface {
	GetUnreadNotifications(userID int) ([]Notification, error)
	MarkNotificationAsRead(notificationID int) error
}

// NotificationService is the concrete implementation of NotificationServe
type NotificationService struct {
	DB *sql.DB
}

// GetUnreadNotifications retrieves unread notifications for a user
func (ns *NotificationService) GetUnreadNotifications(userID int) ([]Notification, error) {
	query := "SELECT id, user_id, message, read FROM notifications WHERE user_id = ? AND read = 0"
	rows, err := ns.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Message, &n.Read); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// MarkNotificationAsRead updates the notification as read
func (ns *NotificationService) MarkNotificationAsRead(notificationID int) error {
	query := "UPDATE notifications SET read = 1 WHERE id = ?"
	_, err := ns.DB.Exec(query, notificationID)
	return err
}
