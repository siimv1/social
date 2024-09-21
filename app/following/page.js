"use client";
import React, { useState, useEffect } from 'react';
import { apiRequest } from '../apiclient';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import './following.css';
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
const FollowPage = () => {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const router = useRouter();
    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const response = await apiRequest("/user", "GET");
                const data = await response; 
                setUsers(data.users);
                setLoading(false);
            } catch (error) {
                console.error("Failed to fetch users:", error);
                setLoading(false);
            }
        };
        fetchUsers();
    }, []);
    const handleBack = () => {
        router.back();
    };
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
            {/* Back Button */}
            <button className="back-button" onClick={handleBack}>
                Back
            </button>
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
            </div>
            {/* Content */}
            <h1 className="followers-header">I'm Following</h1>
            <div className="followers-content">
                    <p>I’m not following anyone yet.</p>
            </div>
        </div>
    );
};
export default FollowPage;
