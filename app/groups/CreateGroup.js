"use client";

import React, { useState } from 'react';
import { apiRequest } from '../apiclient';


const CreateGroup = ({ onGroupCreated }) => {
    const [groupName, setGroupName] = useState('');
    const [groupDescription, setGroupDescription] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
           
            const body = {
                title: groupName,
                description: groupDescription,
            };

          
            const response = await apiRequest('/groups/create', 'POST', body);

            if (response && response.id) {
                onGroupCreated(response); 
                setGroupName(''); 
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

export default CreateGroup; 