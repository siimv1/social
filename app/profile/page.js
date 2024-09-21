"use client";
import Link from 'next/link';
import Image from 'next/image';
import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import './home.css';

const Home = () => {
    const router = useRouter();
    const [profileData, setProfileData] = useState(null);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);
    const [isPublicProfile, setIsPublicProfile] = useState(true);

    const isOwnProfile = true;  

    const toggleProfileVisibility = () => {
        setIsPublicProfile(prev => !prev);
    };

    const handleLogout = async () => {
        router.push('/login');
    };

    const handleBack = () => {
        router.back();
    };

    // Profiili andmete hankimine
    useEffect(() => {
        const fetchProfileData = async () => {
            const token = localStorage.getItem('token');
            if (!token) {
                router.push('/login');
                return;
            }

            try {
                const response = await fetch('http://localhost:8080/profile', {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    }
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    throw new Error(errorData.message || 'Failed to fetch profile');
                }

                const data = await response.json();
                setProfileData(data); // Salvestame andmed profiili jaoks
                setLoading(false);
            } catch (error) {
                setError(error.message);
                setLoading(false);
            }
        };

        fetchProfileData();
    }, [router]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div className="home-container">
            <div className="home-header">
                <Link href="/profile">
                    <Image src="/profile.png" alt="Profile" width={100} height={100} className="profile-pic" />
                </Link>
                <h1>Social Network</h1>
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
                    <h2>{profileData.first_name} {profileData.last_name}</h2>
                    <p>Nickname: {profileData.nickname}</p>
                    <p>Email: {profileData.email}</p>
                    <p>Date of birth: {profileData ? new Date(profileData.date_of_birth).toISOString().split('T')[0] : 'Date of birth'}</p>
                    <p>About Me: {profileData.about_me}</p>
                    {profileData.avatar && (
                        <img src={profileData.avatar} alt="Profile Avatar" className="avatar" />
                    )}
                </div>

                {isOwnProfile && (
                    <div className="profile-settings">
                        <button onClick={toggleProfileVisibility}>
                            {isPublicProfile ? "Switch to Private Profile" : "Switch to Public Profile"}
                        </button>
                    </div>
                )}

                <button type="button" onClick={handleBack} className="back-button" style={{ marginTop: '10px' }}>Back</button>
            </div>

            <div className="home-sidebar-right">
                <div className="followers-following">
                    <h2>Followers</h2>
                    <p>Follower 1</p>
                    <p>Follower 2</p>

                    <h2>Following</h2>
                    <p>Following 1</p>
                    <p>Following 2</p>
                </div>
            </div>

            <div className="home-content">
                <div className="user-posts">
                    <h2>My posts</h2>
                    <div className="post">
                        <p>Post 1</p>
                    </div>
                    <div className="post">
                        <p>Post 2</p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;
