"use client";

import React, { useEffect, useState } from 'react';
import { apiRequest } from '../apiclient';
import Link from 'next/link';
import './request.css';

const PendingFollowRequests = ({ profileUserId }) => {
    const [requests, setRequests] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [loggedInUserId, setLoggedInUserId] = useState(null);

    // Fetch the session to get the logged-in user ID
    useEffect(() => {
        const fetchSession = async () => {
            try {
                const sessionResponse = await apiRequest('/session', 'GET');
                const userId = sessionResponse.user_id;
                setLoggedInUserId(userId);
            } catch (error) {
                console.error('Failed to fetch session:', error);
            }
        };

        fetchSession();
    }, []);

    // Check if the profile being viewed belongs to the logged-in user
    const isOwnProfile = loggedInUserId === profileUserId;

    useEffect(() => {
        if (isOwnProfile) {
            console.log('Fetching pending follow requests...');
            fetchPendingRequests();
        }
    }, [isOwnProfile]);
    const fetchPendingRequests = async () => {
        try {
            const data = await apiRequest('/followers/requests', 'GET');
            console.log('Pending follow requests data:', data);
            setRequests(data.requests || []);
        } catch (error) {
            console.error('Failed to fetch pending follow requests:', error.message);
            setError('Failed to fetch pending follow requests');
        } finally {
            setLoading(false);
        }
    };

    const handleAccept = async (followerId) => {
        try {
            await apiRequest('/followers/accept', 'POST', { follower_id: followerId });
            // Update the requests list
            setRequests(prevRequests => prevRequests.filter(req => req.id !== followerId));
        } catch (error) {
            console.error('Failed to accept follow request:', error.message);
        }
    };

    const handleReject = async (followerId) => {
        try {
            await apiRequest('/followers/reject', 'POST', { follower_id: followerId });
            // Update the requests list
            setRequests(prevRequests => prevRequests.filter(req => req.id !== followerId));
        } catch (error) {
            console.error('Failed to reject follow request:', error.message);
        }
    };

    // Return null if the profile being viewed is not the logged-in user's profile
    if (!isOwnProfile) {
        return null;
    }

    if (loading) {
        return <p>Loading pending follow requests...</p>;
    }

    if (error) {
        return <p>{error}</p>;
    }

    return (
        <div className="pending-requests">
            <h2>Pending Follow Requests</h2>
            {requests.length === 0 ? (
                <p>No pending follow requests.</p>
            ) : (
                requests.map(req => (
                    <div key={req.id} className="request-item">
                        <p>
                            <Link href={`/profile/${req.id}`}>
                                {req.first_name} {req.last_name}
                            </Link>
                        </p>
                        <button className="accept-button" onClick={() => handleAccept(req.id)}>Accept</button>
                        <button className="reject-button" onClick={() => handleReject(req.id)}>Reject</button>
                    </div>
                ))
            )}
        </div>
    );
};

export default PendingFollowRequests;
