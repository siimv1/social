import axios from 'axios';

const NotificationService = () => {
  const getUnreadNotifications = async (userId) => {
    try {
      const response = await axios.get(`/notifications/unread?user_id=${userId}`);
      return response.data;
    } catch (error) {
      console.error("Error fetching unread notifications:", error);
      return [];
    }
  };

  const markNotificationAsRead = async (id) => {
    try {
      await axios.post(`/notifications/read?id=${id}`);
    } catch (error) {
      console.error("Error marking notification as read:", error);
    }
  };

  return {
    getUnreadNotifications,
    markNotificationAsRead,
  };
};

export default NotificationService;
