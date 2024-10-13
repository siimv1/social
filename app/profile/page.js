"use client";

import Link from 'next/link';
import Image from 'next/image';
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation'; 
import './profile.css';
import { apiRequest } from '../apiclient';
import PendingFollowRequests from '../requests/page.js';
import PostList from '../posts/PostList';  

const Home = () => {
    const router = useRouter();
    const [profileData, setProfileData] = useState(null);
    const [followers, setFollowers] = useState([]);
    const [following, setFollowing] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);
    const [isPublicProfile, setIsPublicProfile] = useState(true);

    const toggleProfileVisibility = async () => {
        try {
            const newVisibility = !isPublicProfile;
            const response = await apiRequest('/profile/visibility', 'POST', { isPublic: newVisibility });
            if (response.success) {
                setIsPublicProfile(newVisibility);
            } else {
                console.error('Failed to update profile visibility');
            }
        } catch (error) {
            console.error('Error updating profile visibility:', error.message);
        }
    };

    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    useEffect(() => {
        const fetchProfileData = async () => {
            const token = localStorage.getItem('token');
            if (!token) {
                router.push('/login');
                return;
            }
        
            try {
                const data = await apiRequest('/profile', 'GET');
                setProfileData(data);
                setIsPublicProfile(data.is_public);
                setLoading(false);
            } catch (error) {
                setError(error.message);
                setLoading(false);
            }
        };

        const fetchFollowers = async () => {
            try {
                const data = await apiRequest("/followers/list", "GET");
                setFollowers(data.followers || []);
            } catch (error) {
                console.error("Failed to fetch followers:", error.message);
            }
        };

        const fetchFollowing = async () => {
            try {
                const data = await apiRequest("/following/list", "GET");
                setFollowing(data.following || []);
            } catch (error) {
                console.error("Failed to fetch following:", error.message);
            }
        };

        fetchProfileData();
        fetchFollowers();
        fetchFollowing();
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div className="home-container">
            <div className="home-header">
                <Link href={`/profile/${profileData?.id}`}>
                    <Image src="/profile.png" alt="Profile" width={100} height={100} className="profile-pic" />
                </Link>
                <Link href="/home" style={{ textDecoration: 'none', color: 'inherit' }} >
                    <h1>Social Network</h1>
                </Link>
                <div className="header-buttons">
                    <button className="notification-button">
                        <Image src="/notification.png" alt="Notifications" width={40} height={40} />
                    </button>
                    <button className="logout-button" onClick={handleLogout}>Log Out</button>
                </div>
            </div>

            <div className="home-sidebar-left">
                <div className="profile-info">
                    <h2>{profileData?.first_name} {profileData?.last_name}</h2>
                    <p>Nickname: {profileData?.nickname}</p>
                    <p>Email: {profileData?.email}</p>
                    <p>Date of birth: {profileData?.date_of_birth ? new Date(profileData.date_of_birth).toLocaleDateString() : 'N/A'}</p>
                    <p>About Me: {profileData?.about_me}</p>
                    {profileData?.avatar && (
                        <img src={profileData.avatar} alt="Profile Avatar" className="avatar" />
                    )}

                    {/* Include the PendingFollowRequests component */}
                    <PendingFollowRequests profileUserId={profileData?.id} />
                </div>
                <div className="profile-settings">
                    <button onClick={toggleProfileVisibility}>
                        {isPublicProfile ? 'Switch to Private Profile' : 'Switch to Public Profile'}
                    </button>
                    <button type="button" onClick={() => window.history.back()} className="back-button" style={{ marginTop: '10px' }}>
                        Back
                    </button>
                </div>
            </div>

            <div className="home-sidebar-right">
                <h2>Following</h2>
                {following.length > 0 ? (
                    following.map(followed => (
                        <p key={followed.id}>
                            <Link href={`/profile/${followed.id}`}>
                                {followed.first_name} {followed.last_name}
                            </Link>
                        </p>
                    ))
                ) : <p>Not following anyone yet.</p>}
    
                <h2>Followers</h2>
                {followers.length > 0 ? (
                    followers.map(follower => (
                        <p key={follower.id}>
                            <Link href={`/profile/${follower.id}`}>
                                {follower.first_name} {follower.last_name}
                            </Link>
                        </p>
                    ))
                ) : <p>No followers yet.</p>}
            </div>

            <div className="home-content">
                <div className="user-posts">
                    <h2>My posts</h2>
                    {profileData && <PostList userId={profileData.id} />}
                </div>
            </div>
        </div>
    );
};

export default Home;
