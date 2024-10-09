"use client";

import Link from 'next/link';
import Image from 'next/image';
import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
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
    const [showChat, setShowChat] = useState(false); // Toggle state for chat box
    const router = useRouter();

    const handleMessengerClick = () => {
        setShowChat(!showChat); // Toggle the chat box on and off
    };

    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    const handlePostCreated = (post) => {
        setNewPost(post);
    };

    const handleFollowChange = (userId, isFollowing) => {
        setUsers(prevUsers =>
            prevUsers.map(user =>
                user.id === userId ? { ...user, isFollowing } : user
            )
        );
    };

    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const data = await apiRequest("/users", "GET");
                if (data && data.users) {
                    const usersWithDefaults = data.users.map(user => ({
                        ...user,
                        isFollowing: user.isFollowing ?? false,
                        isOnline: user.isOnline ?? false,
                    }));
                    setUsers(usersWithDefaults);

                    const loggedInUserId = parseInt(localStorage.getItem("userId")); 
                    const currentUser = usersWithDefaults.find(u => u.id === loggedInUserId);
                    if (currentUser) {
                        setUser(currentUser);
                    }
                } else {
                    console.log("No users found");
                }
                setLoading(false);
            } catch (error) {
                console.error("Failed to fetch users:", error);
                setLoading(false);
            }
        };

        fetchUsers();
      

        const fetchFollowers = async () => {
            try {
                const data = await apiRequest("/followers/list", "GET");
                setFollowers(data.followers);
            } catch (error) {
                console.error("Failed to fetch followers:", error.message);
            }
        };

        fetchFollowers();
    }, []);

    const filteredUsers = users.filter(u => u.id !== user?.id);

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
                <button className="logout-button" onClick={handleLogout}>Log Out</button>
            </div>

            <div className="home-sidebar-left">
                <ul>
                    <li><Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }} >My profile</Link></li>
                    <li><Link href="/groups" style={{ textDecoration: 'none', color: 'inherit' }} >Groups</Link></li>
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
                        filteredUsers.map(user => (
                            <div key={user.id} className={`user-item ${user.isOnline ? 'online' : 'offline'}`}>
                                <Link href={`/profile/${user.id}`} style={{ textDecoration: 'none', color: 'inherit' }} >
                                    <p>{user.first_name} {user.last_name}</p>
                                </Link>
                            </div>
                        ))
                    )}
                </div>
            </div>

            <div className="home-content">
                <div className="post-section">
                    <h2>Create a Post</h2>
                    {user ? <CreatePost onPostCreated={handlePostCreated} userId={user.id} /> : <p>Loading user data...</p>}
                </div>

                <div className="timeline-section">
                    <h2>Your Timeline</h2>
                    {user ? (
                        <PostList userId={user.id} newPost={newPost} />
                    ) : (
                        <p>Loading user data...</p>
                    )}
                </div>

                {showChat && (
                    <div className="chat-section">
                        <Chat /> {/* Include the Chat component */}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Home;
