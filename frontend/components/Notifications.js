import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import NotificationService from '../../utilities/notification_service';
import { setNotifications, markAsRead } from '../../store/notificationSlice';

const Notifications = ({ userId }) => {
  const dispatch = useDispatch();
  const notifications = useSelector((state) => state.notifications.notifications);
  const notificationService = NotificationService();

  useEffect(() => {
    const fetchNotifications = async () => {
      const unreadNotifications = await notificationService.getUnreadNotifications(userId);
      dispatch(setNotifications(unreadNotifications));
    };

    fetchNotifications();
  }, [dispatch, userId]);

  const handleMarkAsRead = async (id) => {
    await notificationService.markNotificationAsRead(id);
    dispatch(markAsRead(id));
  };

  return (
    <div>
      <h3>Notifications</h3>
      <ul>
        {notifications.map((notification) => (
          <li key={notification.id}>
            <span>{notification.message}</span>
            {!notification.read && (
              <button onClick={() => handleMarkAsRead(notification.id)}>Mark as Read</button>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Notifications;
