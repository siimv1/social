"use client";


import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '../../apiclient';
import Link from 'next/link';
import Image from 'next/image';
import './groups.css';

const GroupDetail = () => {
    const router = useRouter();
    const [group, setGroup] = React.useState(null);
    const [allGroups, setAllGroups] = useState([]);
    const [loading, setLoading] = React.useState(true);

    const loadGroup = async (id) => {
        if (!id) return;
        try {
            const response = await apiRequest(`/groups/${id}`, 'GET');
            console.log('Group data loaded:', response);
            setGroup(response);
            setLoading(false);
        } catch (error) {
            console.error('Failed to load group:', error);
            setLoading(false);
        }
    };

    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    const handleBack = () => {
        router.back();
    };

    React.useEffect(() => {
        if (router.query && router.query.id) {
            loadGroup(router.query.id);
        }
    }, [router.query]);

    const handleJoinGroup = async () => {
        console.log('Joining group:', group.title);
    };

    const handleAddEvent = async () => {
        console.log('Adding Event:', group.title);
    };

    // if (loading) {
    //     return <p>Loading group details...</p>;
    // }

    // if (!group) {
    //     return <p>No group found.</p>;
    // }



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
                    {/* <li><Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }}>My profile</Link></li>
                       */}
                </ul>
                <button type="button" onClick={handleBack} className="back-button" style={{ marginTop: '10px' }}>
                    Back
                </button>
            </div>

            {/* Right Sidebar */}
            <div className="home-sidebar-right">


                <button onClick={handleJoinGroup} className="join-button">Join Group</button>
                <br>
                </br>
                <button onClick={handleAddEvent} className="event-button" >Add Event</button>
                <br>
                </br>

                <div className="group-section">
                    <h4>All Events</h4>
                </div>
            </div>

            {/* Main Content */}
            <div className="home-content">
                <div className="group-section">
                    {/* <h1>{group.title}</h1>
                        <p>{group.description}</p> */}
                    <h1>Group Title</h1>

                    {/* <button onClick={handleJoinGroup}>Join Group</button> */}
                </div>
                <br></br>
                <div className="group-section">
                    <h4>description</h4>
                </div>
                <br></br>
                <div className="group-section">
                    <h4>posts</h4>
                </div>
            </div>
        </div>
    );
};


export default GroupDetail;