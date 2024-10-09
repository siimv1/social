"use client";

import React, { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { apiRequest } from '../../apiclient'; 
import Link from 'next/link';
import Image from 'next/image';
import './groups.css';

const GroupDetail = () => {
    const router = useRouter();
    const params = useParams(); 
    const { id } = params; 
    const [group, setGroup] = useState(null);
    const [loading, setLoading] = useState(true);

    const loadGroup = async (groupId) => {
        if (!groupId) {
            console.log('No ID provided.');
            return;
        }
        console.log('Fetching group with ID:', groupId);
        try {
            const response = await apiRequest(`/groups/${groupId}`, 'GET');
            console.log('API response:', response);

            if (response && typeof response === 'object' && response.id) {
                setGroup(response); 
            } else {
                console.log('Unexpected response format or empty response:', response);
            }
        } catch (error) {
            console.error('Failed to load group:', error);
        } finally {
            setLoading(false);
        }
    };

    
    useEffect(() => {
        if (id) {
            console.log('Extracted ID:', id);
            loadGroup(id); 
        } else {
            console.log('No ID found in the URL.');
        }
    }, [id]);

    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    const handleBack = () => {
        router.back();
    };

    const handleJoinGroup = async () => {
        console.log('Joining group:', group.title);
        // Implement join group functionality here
    };

    const handleAddEvent = async () => {
        console.log('Adding Event:', group.title);
        // Implement add event functionality here
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
                <button type="button" onClick={handleBack} className="back-button" style={{ marginTop: '10px' }}>
                    Back
                </button>
            </div>

            {/* Right Sidebar */}
            <div className="home-sidebar-right">
                <button onClick={handleJoinGroup} className="join-button">Join Group</button>
                <br />
                <button onClick={handleAddEvent} className="event-button">Add Event</button>
                <br />
                <div className="group-section">
                    <h4>All Events</h4>
                </div>
            </div>

            {/* Main Content */}
            <div className="home-content">
                {group ? (
                    <div className="group-section">
                        <h1>{group.title}</h1>
                        <p>{group.description}</p>
                       
                       
                    </div>
                ) : (
                    <div>No group found.</div> 
                )}
            </div>
        </div>
    );
};

export default GroupDetail;