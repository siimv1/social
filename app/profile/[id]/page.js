"use client";

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter, useParams } from 'next/navigation';
import { apiRequest } from '../../apiclient';
import '../profile.css';
import PostList from '../../posts/PostList.js';
import Chat from '../../chat/Chat';


const UserProfile = () => {
    const router = useRouter(); // Navigatsioonimeetodite jaoks
    const params = useParams(); // Saame route'i parameetrid
    const { id } = params; // Saame ID dünaamilisest route'ist
    const loggedInUserId = 1; 
    const [userData, setUserData] = useState(null);
    const [isFollowing, setIsFollowing] = useState(false); // Jälgimise oleku haldamine
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [followers, setFollowers] = useState([]);
    const [following, setFollowing] = useState([]);
  const [showChat, setShowChat] = useState(false); // State to toggle chat box visibility

    const handleSendMessage = () => {
        setShowChat(true); // Show the chat box when the button is clicked
      };
    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    const handleFollowChange = (followingStatus) => {
        setIsFollowing(followingStatus); // Uuendame jälgimise olekut
    };

    useEffect(() => {
        if (!id) {
            setError('User ID not found');
            setLoading(false);
            return;
        }

        const fetchUserData = async () => {
            try {
                // Toome valitud kasutaja andmed ID järgi
                const data = await apiRequest(`/users/${id}`, 'GET');
                setUserData(data);
                setIsFollowing(data.is_following); // Backend tagastab, kas praegune kasutaja jälgib seda profiili
            } catch (error) {
                setError(error.message);
            } finally {
                setLoading(false);
            }
        };

        fetchUserData();
    }, [id]);

    useEffect(() => {
        if (!id) return;

        const fetchFollowers = async () => {
            try {
                const data = await apiRequest(`/followers/list/${id}`, 'GET');
                setFollowers(data.followers || []);
            } catch (error) {
                console.error("Failed to fetch followers:", error.message);
            }
        };

        fetchFollowers();
    }, [id]);

    useEffect(() => {
        if (!id) return;

        const fetchFollowing = async () => {
            try {
                const data = await apiRequest(`/following/list/${id}`, 'GET');
                setFollowing(data.following || []);
            } catch (error) {
                console.error("Failed to fetch following:", error.message);
            }
        };

        fetchFollowing();
    }, [id]);

    const FollowButton = ({ followedID, isFollowingInitial, onFollowChange }) => {
        const [isFollowingState, setIsFollowingState] = useState(isFollowingInitial);
        const [loadingState, setLoadingState] = useState(false);

        useEffect(() => {
            setIsFollowingState(isFollowingInitial);
        }, [isFollowingInitial]);

        const handleFollow = async () => {
            setLoadingState(true);
            try {
                const response = await apiRequest('/followers', 'POST', { followed_id: followedID });
                if (response.status === 'accepted') {
                    setIsFollowingState(true);
                    onFollowChange(true);
                }
            } catch (error) {
                console.error('Follow request failed:', error.message);
            } finally {
                setLoadingState(false);
            }
        };

        const handleUnfollow = async () => {
            setLoadingState(true);
            try {
                const response = await apiRequest('/followers/unfollow', 'POST', { followed_id: followedID });
                if (response.status === 'OK') {
                    setIsFollowingState(false);
                    onFollowChange(false);
                }
            } catch (error) {
                console.error('Unfollow request failed:', error.message);
            } finally {
                setLoadingState(false);
            }
        };

        return (
            <div>
                {isFollowingState ? (
                    <button className="unfollow-button" onClick={handleUnfollow} disabled={loadingState}>
                        {loadingState ? 'Unfollowing...' : 'Unfollow'}
                    </button>
                ) : (
                    <button className="follow-button" onClick={handleFollow} disabled={loadingState}>
                        {loadingState ? 'Following...' : 'Follow'}
                    </button>
                )}
            </div>
        );
    };

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error || !userData) {
        return <div>Error: {error || 'User data not found'}</div>;
    }

    return (
        <div className="home-container" key={id}>
            <div className="home-header">
                <Link href="/profile">
                    <Image src="/profile.png" alt="Profile" width={100} height={100} className="profile-pic" />
                </Link>
                <Link href="/home" className="home-title-link">
                    <h1>Social Network</h1>
                </Link>
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

            <div className="home-sidebar-left">
                <div className="profile-info">
                    <h2>{userData.first_name} {userData.last_name}</h2>
                    <p>Nickname: {userData.nickname}</p>
                    <p>Email: {userData.email}</p>
                    <p>Date of birth: {new Date(userData.date_of_birth).toISOString().split('T')[0]}</p>
                    <p>About Me: {userData.about_me}</p>
                    {userData.avatar && (
                        <img src={userData.avatar} alt="User Avatar" className="avatar" />
                    )}

                    <FollowButton
                        followedID={parseInt(id)}
                        isFollowingInitial={isFollowing}
                        onFollowChange={handleFollowChange}
                    />
                    
                    <button onClick={handleSendMessage} className="send-message-btn">Send Message</button>

                    {showChat && (
                    <div className="chat-box">
                        {console.log("Rendering Chat Component - Sender ID:", loggedInUserId, "Recipient ID:", id)}
                        <Chat senderId={loggedInUserId} recipientId={id} />
                            </div>
                        )}
                        </div>

                        <button type="button" onClick={() => router.back()} className="back-button" style={{ marginTop: '10px' }}>
                            Back
                            </button>
                    </div>

            <div className="home-sidebar-right">
                <div className="followers-following">
                    <h2>Following</h2>
                    {following.length === 0 ? (
                        <p>Not following anyone yet.</p>
                    ) : (
                        following.map(user => (
                            <p key={user.id}>
                                <Link href={`/profile/${user.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                                    {user.first_name} {user.last_name}
                                </Link>
                            </p>
                        ))
                    )}

                    <h2>Followers</h2>
                    {followers.length === 0 ? (
                        <p>No followers yet.</p>
                    ) : (
                        followers.map(follower => (
                            <p key={follower.id}>
                                <Link href={`/profile/${follower.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                                    {follower.first_name} {follower.last_name}
                                </Link>
                            </p>
                        ))
                    )}
                </div>
            </div>

            <div className="home-content">
                <div className="user-posts">
                    <h2>My posts</h2>
                    {userData && <PostList userId={userData.id} />}
                </div>
            </div>

            {/* Render the chat box dynamically in the bottom-right corner */}
            {showChat && (
                <div className="chat-box">
                    <Chat senderId={loggedInUserId} recipientId={id} />
                </div>
            )}
        </div>
    );
};

export default UserProfile;
