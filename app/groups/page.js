"use client";

import Link from 'next/link';
import Image from 'next/image';
import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '../apiclient';
import CreateGroup from './CreateGroup';
import './groups.css';



const Home = () => {
    const [myGroups, setMyGroups] = useState([]); 
    const [allGroups, setAllGroups] = useState([]); 
    const router = useRouter();

    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    const handleBack = () => {
        router.back();
    };

    const handleGroupCreated = (group) => {
        console.log('New group created:', group);
        setMyGroups(prevGroups => [...prevGroups, group]); 
    };

   
    const loadAllGroups = async () => {
        try {
            const response = await apiRequest('/groups', 'GET');
            setAllGroups(response); 
        } catch (error) {
            console.error('Failed to load groups:', error);
        }
    };

    
    useEffect(() => {
        loadAllGroups();
    }, []);

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
                    <li><Link href="/groups" style={{ textDecoration: 'none', color: 'inherit' }}>Groups</Link></li>
                </ul>
                <button type="button" onClick={handleBack} className="back-button" style={{ marginTop: '10px' }}>
                    Back
                </button>
            </div>

            {/* Right Sidebar */}
            <div className="home-sidebar-right">
                <div>
                    <h2>All Groups</h2>
                    {allGroups.length > 0 ? (
                        <ul>
                            {allGroups.map(group => (
                                <li key={group.id}>{group.title}</li> 
                            ))}
                        </ul>
                    ) : (
                        <p>There are currently no groups available.</p>
                    )}
                </div>
            </div>

            {/* Main Content */}
            <div className="home-content">
                <div className="group-section">
                    <h2>Create a new Group</h2>
                    <CreateGroup onGroupCreated={handleGroupCreated} />
                </div>

                <div className="my-groups">
                    <h2>My Groups</h2>
                    {myGroups.length > 0 ? (
                        <ul>
                            {myGroups.map(group => (
                                <li key={group.id}>{group.title}</li> 
                            ))}
                        </ul>
                    ) : (
                        <p>No groups created yet.</p>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Home;