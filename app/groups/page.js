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
    const [userId, setUserId] = useState(null); 
    const router = useRouter();

    useEffect(() => {
        const checkSession = async () => {
            try {
                const sessionResponse = await apiRequest('/session', 'GET');
                const currentUserId = sessionResponse.user_id;
                if (currentUserId) {
                    setUserId(currentUserId);
                } else {
                    router.push('/login'); 
                }
            } catch (error) {
                console.error('Error fetching session:', error);
                router.push('/login');
            }
        };
        
        checkSession();
    }, []);  // Käivitub ainult kord, kui komponent monteeritakse    
    

    const handleLogout = async () => {
        try {
            await apiRequest('/logout', 'POST');
            setUserId(null);
            router.push('/login');
        } catch (error) {
            console.error('Failed to log out:', error);
        }
    };

    const handleBack = () => {
        router.back();
    };

    const handleGroupCreated = (group) => {
        console.log('New group created:', group);
    
        // Kontrolli, kas grupp on juba myGroups hulgas olemas
        setMyGroups((prevGroups) => {
            const isGroupAlreadyPresent = prevGroups.some(g => g.id === group.id);
            if (!isGroupAlreadyPresent) {
                // Lisa uus grupp ainult siis, kui see pole juba olemas
                return [...prevGroups, group];
            }
            return prevGroups;
        });
    
        // Laadi kõik grupid uuesti
        setTimeout(() => {
            loadAllGroups(); // Edastame loodud grupi asemel lihtsalt kõik grupid
        }, 500);  // 500ms ooteaeg
    };
    

    useEffect(() => {
        if (typeof window !== 'undefined') {  
            const storedUserId = localStorage.getItem('userId');
            if (storedUserId) {
                setUserId(parseInt(storedUserId));
            } else {
                router.push('/login'); 
            }
        }
    }, []);

    const loadAllGroups = async (createdGroup = null) => {
        try {
            const response = await apiRequest('/groups', 'GET');
            if (response) {
                // Kontrolli, et userId on saadaval
                if (!userId) {
                    console.error("User ID is not available");
                    return;
                }
    
                // Filtreeri grupid kasutaja ID alusel
                const myGroups = response.filter(group => group.creator_id === userId); 
                const otherGroups = response.filter(group => group.creator_id !== userId);
    
                // Kui uus grupp on loodud, lisame selle myGroups hulka
                if (createdGroup && createdGroup.creator_id === userId) {
                    myGroups.push(createdGroup);
                }
    
                setMyGroups(myGroups);  // Seadista "My Groups"
                setAllGroups(otherGroups);  // Seadista "Other Groups"
            }
        } catch (error) {
            console.error('Failed to load groups:', error);
        }
    };    
    

    useEffect(() => {
        if (userId) {
            loadAllGroups();  // Lae grupid, kui kasutaja ID on saadaval
        }
    }, [userId]);  // Käivitu ainult siis, kui userId muutub
    

    // Load all groups when the component mounts
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
                    <h2>All Other Groups</h2>
                    {allGroups.length > 0 ? (
    <ul>
        {allGroups.map(group => (
            <li key={`allGroup-${group.id}`}>
                <Link href={`/groups/${group.id}`} className="home-sidebar-list" style={{ textDecoration: 'none', color: 'inherit' }}>{group.title}</Link> 
            </li>
        ))}
    </ul>
) : (
    <p>No other groups available.</p>
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
            <li key={`myGroup-${group.id}`}>
                <Link href={`/groups/${group.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>{group.title}</Link> 
            </li>
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