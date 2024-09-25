"use client"; // "use client" peab olema faili alguses

import Link from 'next/link';
import Image from 'next/image';
import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '../apiclient'; // Hoidke ainult üks imporditud apiRequest
import CreatePost from '../app/posts'; // Assuming you have the form in app/posts/page.js

import './home.css';


const HomePage = () => {
    const [posts, setPosts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        // Fetch posts from the backend
        const fetchPosts = async () => {
            try {
                const response = await axios.get('/api/posts'); // Adjust endpoint as needed
                setPosts(response.data);
            } catch (error) {
                console.error('Error fetching posts:', error);
                setError('Failed to load posts. Please try again.');
            } finally {
                setLoading(false);
            }
        };
        fetchPosts();
    }, []);
}
const FollowButton = ({ followedID, isFollowing, onFollowChange }) => {
    const [followStatus, setFollowStatus] = useState(isFollowing ? 'following' : 'not-following');

    const handleFollow = async () => {
        try {
            const response = await apiRequest('/followers', 'POST', { followed_id: followedID });
            if (response.status === 'accepted') {
                setFollowStatus('following');
                onFollowChange(followedID, true);
            } else {
                setFollowStatus('pending');
            }
        } catch (error) {
            console.error('Failed to send follow request:', error);
        }
    };

    const handleUnfollow = async () => {
        try {
            const response = await apiRequest('/followers/unfollow', 'POST', { followed_id: followedID });
            setFollowStatus('not-following');
            onFollowChange(followedID, false);
        } catch (error) {
            console.error('Failed to unfollow:', error);
        }
    };

    return (
        <div>
            {followStatus === 'following' ? (
                <button className="unfollow-button" onClick={handleUnfollow}>Unfollow</button>
            ) : followStatus === 'pending' ? (
                <button className="pending-button" disabled>Pending Request</button>
            ) : (
                <button className="follow-button" onClick={handleFollow}>Follow</button>
            )}
        </div>
    );
};



    useEffect(() => {
        console.log(users.map(user => ({ id: user.id, isOnline: user.isOnline })));
    }, [users]);

    const handleFollowChange = (userId, isFollowing) => {
        setUsers(prevUsers => 
            prevUsers.map(user => 
                user.id === userId ? { ...user, isFollowing, isOnline: user.isOnline ?? false } : user
            )
        );
    };

    return (
        <div className="home-container">
            {/* Header */}
            <div className="home-header">
                <Link href="/profile">
                    <Image src="/profile.png" alt="profile" width={100} height={100} className="profile-pic" />
                </Link>
                <a href="/home" style={{ textDecoration: 'none', color: 'inherit' }} className="home-title-link">
                    <h1>Social Network</h1>
                </a>
                <div className="header-buttons">
                    <button className="notification-button">
                        <Image src="/notification.png" alt="Notifications" width={40} height={40} />
                    </button>
                    <button className="messenger-button">
                        <Image src="/messenger.png" alt="Messenger" width={50} height={50} />
                    </button>
                </div>
                <button className="logout-button" onClick={handleLogout}>Log Out</button>
            </div>

            {/* Left Sidebar */}
            <div className="home-sidebar-left">
                <ul>
                    <li><Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }}>My profile</Link></li>
                    <li><Link href="/followers" style={{ textDecoration: 'none', color: 'inherit' }}>My followers</Link></li> 
                    <li><Link href="/following" style={{ textDecoration: 'none', color: 'inherit' }}>I'm following</Link></li> 
                </ul>
            </div>

            {/* Right Sidebar */}
            <div className="home-sidebar-right">
                <ul>
                    <li>Groups</li>
                    <li>Chats</li>
                </ul>
                <div>
                    <h2>Users</h2>
                    {loading ? (
                        <p>Loading users...</p>
                    ) : users.length === 0 ? (
                        <p>There are currently no registered users.</p>
                    ) : (
                        users.map(user => (
                            <div key={user.id} className={`user-item ${user.isOnline ? 'online' : 'offline'}`}>
                                <p>{user.first_name} {user.last_name}</p>
                                <FollowButton followedID={user.id} isFollowing={user.isFollowing} onFollowChange={handleFollowChange} />
                            </div>
                        ))                                             
                    )}
                </div>
            </div>

            {/* Main Content */}
            <div className="home-content">
                <div className="post-section">
                    <h2>Create a Post</h2>
                    
                </div>

                <div className="timeline-section">
                    <h2>Your Timeline</h2>
                    {/* Removed the hardcoded "John Doe" */}
                </div>
            </div>
        </div>
    );


export default Home;