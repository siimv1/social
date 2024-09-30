"use client";
import Link from 'next/link';
import Image from 'next/image';
import React from 'react';
import { useRouter } from 'next/navigation';
import '../app/global.css'; // Adjust the path if necessary

const Header = ({ profileData, handleLogout }) => {
    const router = useRouter();

    return (
        <div className="home-header">
            <Link href="/profile">
                <Image
                    src={profileData.avatar || '/profile.png'}
                    alt="Profile"
                    width={100}
                    height={100}
                    className="profile-pic"
                />
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
            <button className="logout-button" onClick={handleLogout}>
                Log Out
            </button>
        </div>
    );
};

export default Header;
