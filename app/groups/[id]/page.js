"use client";
import React, { useState, useEffect, useRef } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { apiRequest } from '../../apiclient';
import Link from 'next/link';
import Image from 'next/image';
import './groups.css';
import CreatePost from '../../posts/CreatePost';  // Import CreatePost component
import PostList from '../../posts/PostList';      // Import PostList component

const GroupDetail = () => {
    const router = useRouter();
    const params = useParams();
    const { id } = params;
    const [group, setGroup] = useState(null);
    const [userId, setUserId] = useState(null);
    const [loading, setLoading] = useState(true);
    const [users, setUsers] = useState([]);
    const [selectedUser, setSelectedUser] = useState("");
    const [joinRequested, setJoinRequested] = useState(false);
    const [requests, setRequests] = useState([]);
    const [members, setMembers] = useState([]);
    const [newPost, setNewPost] = useState(null);  // New state to handle newly created posts
    const isCreator = group && group.creator_id === userId;
    const [isMember, setIsMember] = useState(false);
    const userIdLoaded = useRef(false);

    const loadGroup = async (groupId) => {
        if (!groupId) {
            console.log('No ID provided.');
            return;
        }
        try {
            const response = await apiRequest(`/groups/${groupId}`, 'GET');
            if (response && response.id) {
                setGroup(response);
            }
        } catch (error) {
            console.error('Failed to load group:', error);
        } finally {
            setLoading(false);
        }
    };

    const loadAllUsers = async () => {
        try {
            const response = await apiRequest('/users', 'GET');
            console.log('API Response (raw):', response);
            if (response && Array.isArray(response.users)) {
                console.log('Loaded users:', response.users);
                setUsers(response.users);
            } else {
                console.error('Failed to load users: Response is not an array');
            }
        } catch (error) {
            console.error('Error fetching users:', error);
        }
    };

    const loadGroupMembers = async (groupId) => {
        try {
            const response = await fetch(`/groups/${groupId}/members`, {
                method: 'GET',
                headers: {
                    'Authorization': token,
                },
            });
            const membersData = await response.json();
            setMembers(membersData);

            const memberIds = membersData.map(member => member.id);
            if (memberIds.includes(userId)) {
                setIsMember(true);
            }
        } catch (error) {
            console.error('Failed to load members:', error);
        }
    };

    useEffect(() => {
        const storedUserId = localStorage.getItem('userId');
        if (storedUserId) {
            setUserId(parseInt(storedUserId));
        }
    }, []);

    useEffect(() => {
        if (id) {
            loadGroup(id);
            loadAllUsers();
        }
    }, [id]);

    useEffect(() => {
        if (isCreator) {
            loadJoinRequests();
        }
    }, [isCreator]);

    const handleInviteUser = async () => {
        try {
            const response = await apiRequest(`/groups/${group.id}/invite`, 'POST', { userId: selectedUser });
            console.log('User invited:', response);
        } catch (error) {
            console.error('Failed to invite user:', error);
        }
    };

    const handleLogout = async () => {
        localStorage.removeItem('token');
        router.push('/login');
    };

    const handleBack = () => {
        router.back();
    };

    const handleJoinGroup = async () => {
        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('No token found');
            }

            const headers = {
                'Content-Type': 'application/json',
                'Authorization': token,
            };

            const response = await fetch(`http://localhost:8080/groups/${group.id}/join-request`, {
                method: 'POST',
                headers: headers,
            });
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            console.log('Join request sent:', response);

            localStorage.setItem(`joinRequested_${group.id}`, 'true');
            setJoinRequested(true);
        } catch (error) {
            console.error('Failed to send join request:', error);
        }
    };

    const loadJoinRequests = async () => {
        const storedUserId = localStorage.getItem('userId');
        if (storedUserId && parseInt(storedUserId) !== userId) {
            setUserId(parseInt(storedUserId));
        }

        const savedJoinRequest = localStorage.getItem(`joinRequested_${group.id}`);
        if (savedJoinRequest === 'true') {
            setJoinRequested(true);
        }

        try {
            const token = localStorage.getItem('token');
            if (!token) {
                throw new Error('No token found');
            }

            const headers = {
                'Content-Type': 'application/json',
                'Authorization': token,
            };

            const response = await fetch(`http://localhost:8080/groups/${group.id}/join-requests`, {
                method: 'GET',
                headers: headers,
            });
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}, Message: ${await response.text()}`);
            }

            const data = await response.json();
            setRequests(data);
        } catch (error) {
            console.error('Failed to load join requests:', error);
        }

        if (!group || !group.id) {
            console.error('Group ID is missing or invalid');
            return;
        }
    };

    useEffect(() => {
        if (!userIdLoaded.current) {
            const storedUserId = localStorage.getItem('userId');
            if (storedUserId) {
                setUserId(parseInt(storedUserId));
            }
            userIdLoaded.current = true;
        }

        if (group) {
            const savedJoinRequest = localStorage.getItem(`joinRequested_${group.id}`);
            if (savedJoinRequest === 'true' && !joinRequested) {
                setJoinRequested(true);
            }

            loadJoinRequests();
        }
    }, [group]);

    const handleAcceptRequest = async (groupId, userId) => {
        try {
            await apiRequest(`/groups/${groupId}/accept`, 'POST', { userId });
            loadJoinRequests();
        } catch (error) {
            console.error('Failed to accept request:', error);
        }
    };

    const handleDenyRequest = async (groupId, userId) => {
        try {
            await apiRequest(`/groups/${groupId}/deny`, 'POST', { userId });
            loadJoinRequests();
        } catch (error) {
            console.error('Failed to deny request:', error);
        }
    };

    useEffect(() => {
        if (isCreator) {
            loadJoinRequests();
        }
    }, [isCreator]);

    // Function to handle new post creation
    const handlePostCreated = (post) => {
        setNewPost(post);  // Update state when a new post is created
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
                {!isCreator && !joinRequested && (
                    <button onClick={handleJoinGroup} className="join-button">
                        Join Group
                    </button>
                )}
                {!isCreator && joinRequested && (
                    <span>Request Sent</span>
                )}

                {/* Invite Friends section for group creator */}
                {group && userId && isCreator && (
                    <div>
                        <h3>Invite Friends</h3>
                        <select onChange={(e) => setSelectedUser(e.target.value)} value={selectedUser}>
                            <option value="" disabled>Select a user</option>
                            {Array.isArray(users) && users.length > 0 ? (
                                users
                                    .filter(user => user.id !== userId)
                                    .map(user => (
                                        <option key={user.id} value={user.id}>
                                            {user.first_name} {user.last_name}
                                        </option>
                                    ))
                            ) : (
                                <option value="" disabled>No users available</option>
                            )}
                        </select>
                        <button onClick={handleInviteUser}>Invite</button>
                    </div>
                )}
                {/* Pending Join Requests for group creator */}
                {isCreator && (
                    <div className="request-list">
                        <h3>Pending Join Requests</h3>
                        {requests && requests.length > 0 ? (
                            requests.map((req, index) => {
                                const user = users.find(u => u.id === req.user_id);
                                const key = req.user_id && req.group_id ? `${req.group_id}-${req.user_id}` : `undefined-${index}`;
                                return (
                                    <div key={key} className="request-item">
                                        <p>{user ? `${user.first_name} ${user.last_name}` : 'Unknown user'} wants to join group {req.group_id}</p>
                                        <button onClick={() => handleAcceptRequest(req.group_id, req.user_id)}>Accept</button>
                                        <button onClick={() => handleDenyRequest(req.group_id, req.user_id)}>Deny</button>
                                    </div>
                                );
                            })
                        ) : (
                            <p>No pending requests</p>
                        )}
                    </div>
                )}
            </div>

            {/* Main Content */}
            <div className="home-content">
                {group ? (
                    <div className="group-section">
                        <h1>{group.title}</h1>
                        <p>{group.description}</p>

                        {/* Add CreatePost component for group posts */}
                        <CreatePost onPostCreated={handlePostCreated} userId={userId} groupId={id} />

                        {/* Add PostList component for displaying group posts */}
                        <PostList userId={userId} groupId={id} newPost={newPost} />

                    </div>
                ) : (
                    <div>No group found.</div>
                )}
            </div>
        </div >
    );
};

export default GroupDetail;
