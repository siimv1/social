"use client";

import Link from 'next/link';
import Image from 'next/image';
import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '../apiclient';
import './groups.css';


const CreateGroup = ({ onGroupCreated }) => {
    const [groupName, setGroupName] = useState('');
    const [groupDescription, setGroupDescription] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await apiRequest('/groups/create', 'POST', {
                name: groupName,
                description: groupDescription,
            });
            if (response.status === 'success') {
                onGroupCreated(response.group);
                setGroupName(''); // Tühjendame sisendväljad pärast grupi loomist
                setGroupDescription('');
            }
        } catch (error) {
            console.error('Group creation failed:', error);
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <input
                type="text"
                value={groupName}
                onChange={(e) => setGroupName(e.target.value)}
                placeholder="Group Name"
                required
            />
            <textarea
                value={groupDescription}
                onChange={(e) => setGroupDescription(e.target.value)}
                placeholder="Group Description"
                required
            />
            <button type="submit">Create Group</button>
        </form>
    );
};




const Home = () => {

    
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
    };

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
                    <li><Link href="/profile"style={{ textDecoration: 'none', color: 'inherit' }}>My profile</Link></li>
                    <li><Link href="/groups"style={{ textDecoration: 'none', color: 'inherit' }}>Groups</Link></li>
                </ul>



                <button type="button" onClick={handleBack} className="back-button" style={{ marginTop: '10px' }}>
                    Back
                </button>

            </div>



            {/* Right Sidebar */}
            <div className="home-sidebar-right">
                <div>
                    <h2>All Groups</h2>
                   
                        <p>There are currently no created groups.</p>
                
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
                </div>
                <div className="all-groups">
                    <h2>All Groups</h2>
                </div>
            </div>
           
        </div>
    );
};

export default Home;
