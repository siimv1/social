"use client";
import React, { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { apiRequest } from '../../apiclient';
import Link from 'next/link';
import Image from 'next/image';
import './groups.css';
import CreatePost from '../../posts/CreatePost';
import PostList from '../../posts/PostList';
import CreateEvent from './CreateEvent';
import Chat from '../../chat/Chat'; // Import the Chat component (adjust path as necessary)

const GroupDetail = () => {
    const router = useRouter();
    const params = useParams();
    const { id } = params;
    const [group, setGroup] = useState(null);
    const [userId, setUserId] = useState(null);
    const [users, setUsers] = useState([]);
    const [selectedUser, setSelectedUser] = useState("");
    const [joinRequested, setJoinRequested] = useState(false);
    const [requests, setRequests] = useState([]);
    const [members, setMembers] = useState([]);
    const [newPost, setNewPost] = useState(null);
    const [isCreator, setIsCreator] = useState(false);
    const [isMember, setIsMember] = useState(false);
    const [inviteSent, setInviteSent] = useState(false);
    const [eventInvites, setEventInvites] = useState([]);
    const [showChat, setShowChat] = useState(false); // State to toggle chatbox visibility

    useEffect(() => {
        const loadSession = async () => {
            try {
                const response = await apiRequest('/session', 'GET');
                if (response && response.user_id) {
                    localStorage.setItem('userId', response.user_id);
                    setUserId(response.user_id);
                } else {
                    console.error('Failed to get user ID from session');
                }
            } catch (error) {
                console.error('Error loading session:', error);
            }
        };
        loadSession();
    }, []);

    useEffect(() => {
        if (id && userId) {
            loadGroup(id);
            loadAllUsers();
            loadGroupMembers(id);
            checkJoinRequestStatus();
            loadEventInvites();
        }
    }, [id, userId]);

    useEffect(() => {
        if (id && userId) {
            loadGroup(id);
        }
        setUserId(localStorage.getItem('userId'));
    }, [id]);

    useEffect(() => {
        if (group && userId) {
            setIsCreator(group.creator_id === userId);
            const memberIds = members.map(member => member.id);
            setIsMember(memberIds.includes(userId));
        }
    }, [group, userId, members]);

    const loadGroup = async (groupId) => {
        try {
            const response = await apiRequest(`/groups/${groupId}`, 'GET');
            if (response) {
                setGroup(response);
                setIsCreator(response.creator_id === userId);
            } else {
                console.error('Group not found');
            }
        } catch (error) {
            console.error('Failed to load group:', error);
        }
    };

    const loadAllUsers = async () => {
        try {
            const response = await apiRequest('/users', 'GET');
            if (response && Array.isArray(response.users)) {
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
            const response = await apiRequest(`/groups/${groupId}/members`, 'GET');
            const membersData = await response.json();
            setMembers(membersData);
            const memberIds = membersData.map(member => member.id);
            setIsMember(memberIds.includes(userId));
        } catch (error) {
            console.error('Failed to load members:', error);
        }
    };

    const checkJoinRequestStatus = async () => {
        try {
            const response = await apiRequest(`/groups/${id}/join-status`, 'GET');
            if (response && response.requestPending) {
                setJoinRequested(true);
            }
        } catch (error) {
            console.error('Failed to check join request status:', error);
        }
    };

    const handleInviteUser = async () => {
        try {
            const response = await apiRequest(`/groups/${group.id}/invite`, 'POST', { userId: selectedUser });
            console.log('User invited:', response);
            loadGroupMembers(group.id);
            setInviteSent(true);
        } catch (error) {
            console.error('Failed to invite user:', error);
        }
    };

    const handleLogout = async () => {
        await apiRequest('/logout', 'POST');
        router.push('/login');
    };

    const handleBack = () => {
        router.back();
    };

    const handleJoinGroup = async () => {
        if (joinRequested) return;

        try {
            const response = await apiRequest(`/groups/${group.id}/join-request`, 'POST');
            if (response && response.requestPending) {
                setJoinRequested(true);
            }
        } catch (error) {
            console.error('Failed to send join request:', error.message);
        }
    };

    const loadJoinRequests = async () => {
        try {
            const response = await apiRequest(`/groups/${group.id}/join-requests`, 'GET');
            console.log('Join Requests Response:', response);
            if (response && Array.isArray(response)) {
                setRequests(response);  // Handle the array of join requests
            } else {
                console.error('Unexpected response format:', response);
            }
        } catch (error) {
            console.error('Failed to load join requests:', error);
        }
    };

    useEffect(() => {
        if (isCreator) {
            loadJoinRequests();
        }
    }, [isCreator]);

    const handleAcceptRequest = async (userId) => {
        try {
            const response = await apiRequest(`/groups/${group.id}/join-requests/${userId}/accept`, 'POST', { userId });
            if (!response) {
                console.error('Response is null or undefined:', response);
                return;
            }
            if (response.message) {
                loadJoinRequests();
            } else {
                console.error('Unexpected response format:', response);
            }
        } catch (error) {
            console.error('Failed to accept request:', error);
        }
    };

    const handleDenyRequest = async (userId) => {
        try {
            const response = await apiRequest(`/groups/${group.id}/join-requests/${userId}/deny`, 'POST', { userId });
            if (response && response.message) {
                loadJoinRequests();
            } else {
                console.error('Unexpected response format:', response);
            }
        } catch (error) {
            console.error('Failed to deny request:', error);
        }
    };

    const loadEventInvites = async () => {
        try {
            const response = await apiRequest(`/eventinvites`, 'GET');
            if (response && response.eventinvites) {
                setEventInvites(response.eventinvites);
            }
        } catch (error) {
            console.error('Failed to load event invites:', error);
        }
    };

    const handleEventCreated = (eventId) => {
        setShowCreateEvent(false);
    };

    const handleAcceptEventInvite = async (inviteId) => {
        try {
            await apiRequest(`/eventinvites/${inviteId}/accept`, 'POST');
            setEventInvites(eventInvites.filter(inv => inv.id !== inviteId));
        } catch (error) {
            console.error('Failed to accept event invite:', error);
        }
    };

    const handleDeclineEventInvite = async (inviteId) => {
        try {
            await apiRequest(`/eventinvites/${inviteId}/decline`, 'POST');
            setEventInvites(eventInvites.filter(inv => inv.id !== inviteId));
        } catch (error) {
            console.error('Failed to decline event invite:', error);
        }
    };

    // Toggle chatbox visibility
    const handleToggleChat = () => {
        setShowChat((prevShowChat) => !prevShowChat); // Toggle the chatbox visibility
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

            {/* Invite Message */}
            {inviteSent && (
                <div className="invite-message">
                    Kutse on saadetud!
                </div>
            )}

            {/* Left Sidebar */}
            <div className="home-sidebar-left">
                <button type="button" onClick={handleBack} className="back-button" style={{ marginTop: '10px' }}>
                    Back
                </button>
            </div>

            {/* Right Sidebar */}
            <div className="home-sidebar-right">
                {/* Join Group button only if user is not the creator and not a member */}
                {!isCreator && !isMember && !joinRequested && (
                    <button onClick={handleJoinGroup} className="join-button">
                        Join Group
                    </button>
                )}
                {!isCreator && !isMember && joinRequested && (
                    <span>Request Sent</span>
                )}

                {/* Invite Friends section for group creator */}
                {group && userId && isCreator && (
                    <div className="invite-friends">
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
                {isCreator && group && (
                    <div className="request-list">
                        <h3>Pending Join Requests</h3>
                        {requests && requests.length > 0 ? (
                            requests.map((req) => {
                                const user = users.find(u => u.id === req.user_id);
                                return (
                                    <div key={req.user_id} className="request-item">
                                        <p>
                                            {user ? `${user.first_name} ${user.last_name}` : 'Unknown user'} wants to join group {group.title}
                                        </p>
                                        <div className="request-buttons">
                                            <button className="accept" onClick={() => handleAcceptRequest(req.user_id)}>Accept</button>
                                            <button className="deny" onClick={() => handleDenyRequest(req.user_id)}>Deny</button>
                                        </div>
                                    </div>
                                );
                            })
                        ) : (
                            <p>No pending requests</p>
                        )}
                    </div>
                )}

                {/* Events */}
                {(isCreator || isMember) && group && (
                    <CreateEvent
                        groupId={group.id}
                        userId={userId}
                        onEventCreated={handleEventCreated}
                    />
                )}

                {eventInvites && eventInvites.length > 0 && (
                    <div className="event-invites">
                        <h3>Event Invites</h3>
                        {eventInvites.map((invite) => (
                            <div key={invite.id} className="event-invite-item">
                                <p>You are invited to event: {invite.event_title}</p>
                                <button onClick={() => handleAcceptEventInvite(invite.id)}>Accept</button>
                                <button onClick={() => handleDeclineEventInvite(invite.id)}>Decline</button>
                            </div>
                        ))}
                    </div>
                )}

                {/* Group Chat Button */}
                <button onClick={handleToggleChat} className="group-chat-button">
                    {showChat ? 'Hide Group Chat' : 'Join Group Chat'}
                </button>

                {/* Chatbox: Conditionally render chatbox */}
                {showChat && (
    <div className="chatbox-group"> {/* chatbox-group class */}
        <Chat senderId={userId} groupId={group.id} isGroupChat={true} />
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
                        <CreatePost onPostCreated={setNewPost} userId={userId} groupId={id} />

                        {/* Add PostList component for displaying group posts */}
                        <PostList userId={userId} groupId={id} newPost={newPost} />

                    </div>
                ) : (
                    <div>No group found.</div>
                )}
            </div>
        </div>
    );
};

export default GroupDetail;
