"use client";
import Link from 'next/link';
import Image from 'next/image';
import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import './home.css';

const Home = () => {
    const router = useRouter();

    const [isPublicProfile, setIsPublicProfile] = useState(true);
    const isOwnProfile = true; 



    const handleLogout = async () => {
        router.push('/login'); 
    };

    const toggleProfileVisibility = () => {
        setIsPublicProfile(prev => !prev); // Kasutage eelmist väärtust
    };

    const handleBack = () => {
        router.back(); 
    };

    return (
        <div className="home-container">
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

    <div className="home-sidebar-left">
    <div className="profile-info">
            <h2>Name</h2>
            <p>Nickname</p>
            <p>Email</p>
            <p>Date of birth</p>
            <p>About me</p>
        </div>

        <div className="user-activity">
            
            <p>Joined: January 2024</p>
            <p>Last Login: September 2024</p>

              
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
