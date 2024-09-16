"use client";
import Link from 'next/link';
import Image from 'next/image';
import React from 'react';
import { useRouter } from 'next/navigation';
import './home.css';

const Home = () => {
    const router = useRouter();

    const handleLogout = async () => {
        router.push('/login'); 
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
                <ul>
                <li><Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }}>My profile</Link></li>
                    <li>I'm following</li>
                    <li>My followers</li>
                </ul>
            </div>

            <div className="home-sidebar-right">
                <ul>
                    <li>Groups</li>
                    <li>Chats</li>
                </ul>
            </div>

            <div className="home-content">
                <div className="post-section">
                    <h2>Create a Post</h2>
                    <textarea placeholder="What's on your mind?" rows="3"></textarea>
                    <button className="post-button">Post</button>
                </div>

                <div className="timeline-section">
                    <h2>Your Timeline</h2>
                    <div className="post">
                        <h3>John Doe</h3>
                        <p>This is a sample post in your timeline. Like, comment, or share!</p>
                    </div>
                    <div className="post">
                        <h3>Jane Smith</h3>
                        <p>Another sample post with a similar style to Facebook's feed.</p>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Home;
