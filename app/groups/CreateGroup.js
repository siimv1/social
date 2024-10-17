"use client";

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiRequest } from '../apiclient';

const CreateGroup = ({ onGroupCreated }) => {
    const router = useRouter();

    const [groupName, setGroupName] = useState('');
    const [groupDescription, setGroupDescription] = useState('');
    const [userId, setUserId] = useState(null);
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);

    useEffect(() => {
        const checkSession = async () => {
            try {
                const sessionResponse = await apiRequest('/session', 'GET');
                console.log("Session response:", sessionResponse);  // Lisa log, et kontrollida vastust
                const userId = sessionResponse?.user_id;
                if (!userId || isNaN(userId)) {
                    throw new Error('No active session or invalid user ID');
                }
                setUserId(userId);
            } catch (error) {
                console.error('Error checking session:', error);
                router.push('/login'); // Redirect to login if no session
            }
        };        

        checkSession();
    }, [router]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);

        if (!userId || isNaN(userId)) {
            setError('User ID is missing or invalid.');
            return;
        }

        try {
            const body = {
                title: groupName,
                description: groupDescription,
                creator_id: userId, // Veendu, et userId on kehtiv
            };

            const response = await fetch('http://localhost:8080/groups/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(body),
                credentials: 'include' // Ensure cookies are included with the request
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorText}`);
            }

            const data = await response.json();
            console.log("API Response:", data);

            if (data && data.id) {
                onGroupCreated(data);
                setGroupName('');
                setGroupDescription('');
                setSuccess('Group created successfully!');
            } else {
                setError('Failed to create group.');
            }
        } catch (error) {
            console.error('Group creation failed:', error);
            setError('Group creation failed. Please try again.');
        }
    };
    return (
        <div className="create-group-container">
            <form onSubmit={handleSubmit}>
                <input
                    type="text"
                    value={groupName}
                    onChange={(e) => setGroupName(e.target.value)}
                    placeholder="Group Name"
                    required
                    className="group-input"
                />
                <textarea
                    value={groupDescription}
                    onChange={(e) => setGroupDescription(e.target.value)}
                    placeholder="Group Description"
                    required
                    className="group-textarea"
                />
                <button type="submit" className="create-button">Create Group</button>
            </form>

            {error && <p className="error-message">{error}</p>}
            {success && <p className="success-message">{success}</p>}
        </div>
    );
};

export default CreateGroup;
