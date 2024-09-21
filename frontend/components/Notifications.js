import React, { useEffect, useState } from "react";
import axios from "axios";

function Notifications({ userID }) {
  const [notifications, setNotifications] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchNotifications() {
      try {
        const response = await axios.get(`/notifications/unread?user_id=${userID}`);
        setNotifications(response.data);
      } catch (error) {
        console.error("Error fetching notifications:", error);
      } finally {
        setLoading(false);
      }
    }

    fetchNotifications();
  }, [userID]);

  const markAsRead = async (notificationID) => {
    try {
      await axios.post(`/notifications/read/${notificationID}`);
      setNotifications(notifications.filter(n => n.id !== notificationID));
    } catch (error) {
      console.error("Error marking notification as read:", error);
    }
  };

  if (loading) return <div>Loading...</div>;

  return (
    <div className="notification-dropdown">
      <h3>Notifications</h3>
      {notifications.length === 0 ? (
        <p>No new notifications</p>
      ) : (
        notifications.map(notification => (
          <div key={notification.id} className="notification-item">
            <p>{notification.message}</p>
            <button onClick={() => markAsRead(notification.id)}>Mark as Read</button>
          </div>
        ))
      )}
    </div>
  );
}

export default Notifications;
