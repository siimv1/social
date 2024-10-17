"use client";

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { apiRequest } from '../apiclient';
import './home.css';
import CreatePost from '../posts/CreatePost';
import PostList from '../posts/PostList';
import Chat from '../chat/Chat';

const Home = () => {
  const [newPost, setNewPost] = useState(null);
  const [users, setUsers] = useState([]);
  const [followers, setFollowers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState(null);
  const [showChat, setShowChat] = useState(false);
  const router = useRouter();

  const handleMessengerClick = () => {
    setShowChat(!showChat);
  };

  const handleLogout = async () => {
    try {
      await apiRequest('/logout', 'POST');
      setUser(null);
      router.push('/login');
    } catch (error) {
      console.error('Failed to log out:', error);
    }
  };  

  const handlePostCreated = (post) => {
    setNewPost(post);
  };

  const handleFollowChange = (newStatus) => {
    setFollowStatus(newStatus);
    console.log('Follow status updated to:', newStatus);
    // Optionally refresh followers/following lists
};


  useEffect(() => {
    const fetchSessionAndUsers = async () => {
      try {
        // Fetch the current user session
        const sessionResponse = await apiRequest('/session', 'GET');
        const currentUserId = sessionResponse.user_id;

        // Fetch all users
        const usersResponse = await apiRequest('/users', 'GET');
        if (usersResponse && usersResponse.users) {
          const usersWithDefaults = usersResponse.users.map((user) => ({
            ...user,
            isFollowing: user.isFollowing ?? false,
            isOnline: user.isOnline ?? false,
          }));
          setUsers(usersWithDefaults);

          // Find and set the current user
          const currentUser = usersWithDefaults.find((u) => u.id === currentUserId);
          if (currentUser) {
            setUser(currentUser);
          } else {
            // If the current user is not in the users list, set minimal info
            setUser({ id: currentUserId });
          }
        } else {
          console.log('No users found');
        }
        setLoading(false);
      } catch (error) {
        console.error('Failed to fetch session or users:', error);
        setLoading(false);
      }
    };

    fetchSessionAndUsers();
  }, []);

  const filteredUsers = users.filter((u) => u.id !== user?.id);

  return (
    <div className="home-container">
      {/* Header */}
      <div className="home-header">
        <Link href="/profile">
          <Image src="/profile.png" alt="profile" width={100} height={100} className="profile-pic" />
        </Link>
        <Link href="/home" style={{ textDecoration: 'none', color: 'inherit' }} className="home-title-link">
          <h1>Social Network</h1>
        </Link>
        <div className="header-buttons">
          <button className="notification-button">
            <Image src="/notification.png" alt="Notifications" width={40} height={40} />
          </button>
        </div>
        <button className="logout-button" onClick={handleLogout}>
          Log Out
        </button>
      </div>

      <div className="home-sidebar-left">
        <ul>
          <li>
            <Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }}>
              My profile
            </Link>
          </li>
          <li>
            <Link href="/groups" style={{ textDecoration: 'none', color: 'inherit' }}>
              Groups
            </Link>
          </li>
        </ul>
      </div>

      <div className="home-sidebar-right">
        <div>
          <h2>Other users</h2>
          {loading ? (
            <p>Loading users...</p>
          ) : filteredUsers.length === 0 ? (
            <p>There are currently no registered users.</p>
          ) : (
            filteredUsers.map((user) => (
              <div key={user.id} className={`user-item ${user.isOnline ? 'online' : 'offline'}`}>
                <Link href={`/profile/${user.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                  <p>
                    {user.first_name} {user.last_name}
                  </p>
                </Link>
              </div>
            ))
          )}
        </div>
      </div>

      <div className="home-content">
        <div className="post-section">
          <h2>Create a Post</h2>
          {user ? (
            <CreatePost onPostCreated={handlePostCreated} userId={user.id} />
          ) : (
            <p>Loading user data...</p>
          )}
        </div>

        <div className="timeline-section">
          <h2>Your Timeline</h2>
          {user ? <PostList userId={user.id} newPost={newPost} /> : <p>Loading user data...</p>}
        </div>

        {showChat && (
          <div className="chat-section">
            <Chat />
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;
