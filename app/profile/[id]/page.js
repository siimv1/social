"use client";
import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { useRouter, useParams } from 'next/navigation';
import { apiRequest } from '../../apiclient';
import '../profile.css';
import PostList from '../../posts/PostList.js';
import Chat from '../../chat/Chat';
import PendingFollowRequests from '../../requests/page.js';
import { FaPaperPlane, FaArrowLeft } from 'react-icons/fa';

const UserProfile = () => {
    const router = useRouter();
    const params = useParams();
    const { id } = params;
    const profileUserId = parseInt(id); // Profiili ID
    const [loggedInUserId, setLoggedInUserId] = useState(null);
    const [isOwnProfile, setIsOwnProfile] = useState(false);
    const [userData, setUserData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [followers, setFollowers] = useState([]);
    const [following, setFollowing] = useState([]);
    const [showChat, setShowChat] = useState(false);
    const [followStatus, setFollowStatus] = useState('not-following');

    // Kontrolli kasutaja sessiooni ja määra, kas profiil kuulub sisse loginud kasutajale
    useEffect(() => {
        const checkSession = async () => {
            try {
                const sessionResponse = await apiRequest('/session', 'GET');
                const userId = sessionResponse.user_id; // Kontrolli sessiooni põhjal kasutaja ID-d
                if (!userId) {
                    throw new Error('No active session');
                }
                setLoggedInUserId(userId);
                setIsOwnProfile(userId === profileUserId); // Kontrolli, kas profiil on oma
            } catch (error) {
                console.error('Error checking session:', error);
                router.push('/login'); // Kui sessiooni pole, suuna login lehele
            }
        };

        checkSession();
    }, [profileUserId, router]);

    const handleSendMessage = () => {
        setShowChat(true);
    };

    // Kasutaja välja logimine
    const handleLogout = async () => {
        await apiRequest('/logout', 'POST'); // Logi välja serveris sessioonipõhiselt
        router.push('/login');
    };

    // Jälgimise oleku muutmine
    const handleFollowChange = (newStatus) => {
        setFollowStatus(newStatus); // Uuenda follow status
        console.log('Follow status updated to:', newStatus);
    };

    // Lae profiili andmed ja kontrolli follow staatust
useEffect(() => {
    if (!id) {
        setError('User ID not found');
        setLoading(false);
        return;
    }

    const fetchUserData = async () => {
        try {
            const data = await apiRequest(`/users/${id}`, 'GET');
            if (!data) {
                throw new Error('No user data returned');
            }
            console.log('User data:', data);
            setUserData(data);

            // Kontrolli ja seadista follow status vastavalt API andmetele
            setFollowStatus(data.follow_status || 'not-following');
            console.log('Initial follow status:', data.follow_status);
        } catch (error) {
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };

    fetchUserData();
}, [id]);


    useEffect(() => {
        if (!id) return;

        const fetchFollowers = async () => {
            try {
                const data = await apiRequest(`/followers/list/${id}`, 'GET');
                setFollowers(data.followers || []);
            } catch (error) {
                console.error("Failed to fetch followers:", error.message);
            }
        };

        fetchFollowers();
    }, [id]);

    useEffect(() => {
        if (!id) return;

        const fetchFollowing = async () => {
            try {
                const data = await apiRequest(`/following/list/${id}`, 'GET');
                setFollowing(data.following || []);
            } catch (error) {
                console.error("Failed to fetch following:", error.message);
            }
        };

        fetchFollowing();
    }, [id]);

    // Follow nupu komponent, mis kasutab `isPublic`
    const FollowButton = ({ followedID, initialStatus, onFollowChange, isPublic }) => {
        const [loading, setLoading] = useState(false);
        const [followStatusState, setFollowStatusState] = useState(initialStatus);
    
        useEffect(() => {
            setFollowStatusState(initialStatus); // Seadista nupp vastavalt API vastusele
        }, [initialStatus]);
    
        const handleFollow = async () => {
            if (followStatusState === 'accepted' || followStatusState === 'pending') {
                console.log('Already following or follow request is pending.');
                return; // Väldi uuesti päringu tegemist, kui juba jälgitakse või päring on ootel
            }
    
            setLoading(true);
            try {
                const response = await apiRequest('/followers', 'POST', { followed_id: followedID });
                console.log('Follow API response:', response);
    
                if (response.status === 'pending') {
                    setFollowStatusState('pending');
                    onFollowChange('pending');
                } else if (response.status === 'accepted') {
                    setFollowStatusState('accepted');
                    onFollowChange('accepted');
                } else {
                    console.error('Unexpected follow status:', response.status);
                }
            } catch (error) {
                console.error('Follow request failed:', error.message);
            } finally {
                setLoading(false);
            }
        };
    
        const handleUnfollow = async () => {
            setLoading(true);
            try {
                const response = await apiRequest('/followers/unfollow', 'POST', { followed_id: followedID });
                if (response.status === 'OK') {
                    setFollowStatusState('not-following');
                    onFollowChange('not-following');
                }
            } catch (error) {
                console.error('Unfollow request failed:', error.message);
            } finally {
                setLoading(false);
            }
        };
    
        useEffect(() => {
            setFollowStatusState(initialStatus); // Värskenda nupu olek esialgse väärtuse järgi
        }, [initialStatus]);
        
        // Kuvamine vastavalt follow state-ile
        if (followStatusState === 'pending') {
            return <button disabled>Request Pending</button>;
        }
        
        if (followStatusState === 'accepted') {
            return (
                <button className="unfollow-button" onClick={handleUnfollow} disabled={loading}>
                    {loading ? 'Unfollowing...' : 'Unfollow'}
                </button>
            );
        }
        
        if (isPublic || followStatusState === 'not-following') {
            return (
                <button className="follow-button" onClick={handleFollow} disabled={loading}>
                    {loading ? 'Following...' : 'Follow'}
                </button>
            );
        }
    
        return null;
    };    
    

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error || !userData) {
        return <div>Error: {error || 'User data not found'}</div>;
    }

    return (
        <div className="home-container" key={id}>
            <div className="home-header">
                <Link href="/profile">
                    <Image src="/profile.png" alt="profile" width={100} height={100} className="profile-pic" />
                </Link>
                <Link href="/home" className="home-title-link">
                    <h1>Social Network</h1>
                </Link>
                <div className="header-buttons">
                    <button className="notification-button">
                        <Image src="/notification.png" alt="Notifications" width={40} height={40} />
                    </button>
                </div>
                <button className="logout-button" onClick={handleLogout}>Log Out</button>
            </div>
            <div className="home-sidebar-left">
            <div className="profile-info">
    <h2>{userData.first_name} {userData.last_name}</h2>
    {(userData.is_public || followStatus === 'accepted' || isOwnProfile) ? (
        <>
            {userData.nickname && <p>Nickname: {userData.nickname}</p>}
            {userData.email && <p>Email: {userData.email}</p>}
            {userData.date_of_birth && <p>Date of birth: {new Date(userData.date_of_birth).toLocaleDateString()}</p>}
            {userData.about_me && <p>About Me: {userData.about_me}</p>}
            {userData.avatar && <img src={userData.avatar} alt="User Avatar" className="avatar" />}
        </>
    ) : (
        <p>This profile is private.</p>
    )}
    {!isOwnProfile && (
        <FollowButton
            followedID={profileUserId}
            initialStatus={followStatus} // Kasuta `followStatus` API tagastatud väärtust
            onFollowChange={handleFollowChange}
            isPublic={userData.is_public}
        />
                    
                    )}
                    {isOwnProfile && <PendingFollowRequests profileUserId={profileUserId} />}
                    {!isOwnProfile && (
                        <>
                            <br />
                            {(!userData.is_private || followStatus === 'accepted') && (
                                <button onClick={handleSendMessage} className="send-message-btn">
                                    <FaPaperPlane style={{ marginRight: '5px' }} /> Send Message
                                </button>
                            )}
                        </>
                    )}
                </div>
                <button type="button" onClick={() => router.back()} className="back-button" style={{ marginTop: '10px' }}>
                    <FaArrowLeft style={{ marginRight: '5px' }} /> Back
                </button>
            </div>
            <div className="home-sidebar-right">
    <h2>Following</h2>
    {userData && (userData.is_public || followStatus === 'accepted' || isOwnProfile) ? (
        following.length > 0 ? (
            following.map(followed => (
                <p key={followed.id}>
                    <Link href={`/profile/${followed.id}`}>
                        {followed.first_name} {followed.last_name}
                    </Link>
                </p>
            ))
        ) : <p>Not following anyone yet.</p>
    ) : (
        <p>Following is private.</p>
    )}
    
    <h2>Followers</h2>
    {userData && (userData.is_public || followStatus === 'accepted' || isOwnProfile) ? (
        followers.length > 0 ? (
            followers.map(follower => (
                <p key={follower.id}>
                    <Link href={`/profile/${follower.id}`}>
                        {follower.first_name} {follower.last_name}
                    </Link>
                </p>
            ))
        ) : <p>No followers yet.</p>
    ) : (
        <p>Followers are private.</p>
    )}
</div>
           
<div className="home-content">
    <div className="user-posts">
        <h2>{isOwnProfile ? 'My Posts' : `${userData.first_name}'s Posts`}</h2>
        {(userData.is_public || followStatus === 'accepted' || isOwnProfile) ? (
            <PostList userId={userData.id} />
        ) : (
            <p>Posts are private.</p>
        )}
    </div>
</div>
            {showChat && (
                <div className="chat-box">
                    <button className="close-button" onClick={() => setShowChat(false)}>X</button>
                    <Chat senderId={loggedInUserId} recipientId={id} />
                </div>
            )}
        </div>
    );
};

export default UserProfile;
