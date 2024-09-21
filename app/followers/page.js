"use client";
import React, { useState, useEffect } from 'react';
import { apiRequest } from '../apiclient';
import Link from 'next/link';
import Image from 'next/image';
import './followers.css';  // Kasuta sama CSS
const FollowButton = ({ followedID, isFollowing }) => {
    const [followStatus, setFollowStatus] = useState(isFollowing ? 'following' : 'not-following');
    const handleFollow = async () => {
        try {
            const response = await apiRequest('/followers', 'POST', { followed_id: followedID });
            if (response.status === 'accepted') {
                setFollowStatus('following');
            } else {
                setFollowStatus('pending');
            }
        } catch (error) {
            console.error('Failed to send follow request:', error);
        }
    };
    const handleUnfollow = async () => {
        try {
            await apiRequest('/unfollow', 'POST', { followed_id: followedID });
            setFollowStatus('not-following');
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
const FollowPage = ({ isFollowingPage }) => {
    const [users, setUsers] = useState([]);
    const [registeredUsers, setRegisteredUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const url = isFollowingPage ? '/following' : '/followers';
                const response = await apiRequest(url, 'GET');
                setUsers(response.users);
            } catch (error) {
                console.error('Failed to fetch data:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchUsers();
        const fetchRegisteredUsers = async () => {
            try {
                const response = await apiRequest("/users", "GET");
                setRegisteredUsers(response.users);
            } catch (error) {
                console.error("Failed to fetch registered users:", error);
            }
        };
        fetchRegisteredUsers();
    }, [isFollowingPage]);
    if (loading) {
        return <div>Loading...</div>;
    }
    return (
        <div className="followers-container">
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
                <button className="logout-button">Log Out</button>
            </div>
            {/* Sidebarid */}
            <div className="home-sidebar-left">
                <ul>
                    <li><Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }}>My profile</Link></li>
                    <li><Link href="/followers" style={{ textDecoration: 'none', color: 'inherit' }}>My followers</Link></li> 
                    <li><Link href="/following" style={{ textDecoration: 'none', color: 'inherit' }}>I'm following</Link></li> 
                </ul>
            </div>
            <div className="home-sidebar-right">
                <ul>
                    <li>Groups</li>
                    <li>Chats</li>
                </ul>
            </div>
            {/* Sisu */}
            <div className="followers-content">
                <h1 className="followers-header">{isFollowingPage ? "I'm Following" : 'My Followers'}</h1>
                {users.length === 0 ? (
                    <p>{isFollowingPage ? "I’m not following anyone yet." : "No one is following you yet."}</p>
                ) : (
                    users.map(user => (
                        <div key={user.id} className="follower-item">
                            <p>{user.firstName} {user.lastName}</p>
                            <FollowButton followedID={user.id} isFollowing={user.isFollowing} />
                        </div>
                    ))
                )}
            </div>
        </div>
    );
};
export default FollowPage;
