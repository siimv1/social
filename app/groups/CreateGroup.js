"use client";


import React, { useState, useEffect } from 'react';
import { apiRequest } from '../apiclient';


const CreateGroup = ({ onGroupCreated }) => {

    const [groupName, setGroupName] = useState('');
    const [groupDescription, setGroupDescription] = useState('');
    const [userId, setUserId] = useState(null);  
    const [error, setError] = useState(null);
    const [success, setSuccess] = useState(null);


    // console.log(localStorage.getItem('userId'));
    // console.log("CreateGroup component received userId:", userId);
    // console.log(localStorage.getItem('token'));


    // Hangi kasutaja ID localStorage-st
    useEffect(() => {
        if (typeof window !== 'undefined') {  
            const storedUserId = localStorage.getItem('userId');
            console.log("Stored userId:", storedUserId);

            if (storedUserId) {
                setUserId(parseInt(storedUserId)); 
            } else {
                setError('User ID not found');
                window.location.href = '/login'; 
            }
        }
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);
    
        
        let token = localStorage.getItem('token');
    
        if (!token) {
            setError('Token is missing.');
            return;
        }
    
        if (!userId) {
            setError('User ID is missing.');
            return;
        }
    
        try {
            const body = {
                title: groupName,
                description: groupDescription,
                creator_id: userId,
            };
    
         
            const headers = {
                'Content-Type': 'application/json',
                'Authorization': token, 
            };
    
           
            const response = await fetch('http://localhost:8080/groups/create', {
                method: 'POST',
                headers: headers,
                body: JSON.stringify(body)
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