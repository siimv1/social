"use client";
import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import './register.css';

export const apiRequest = async (endpoint, method, body) => {
    const response = await fetch(`http://localhost:8080${endpoint}`, {
        method,
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Network response was not ok');
    }

    return response.json();
};

const Register = () => {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        firstName: '',
        lastName: '',
        dateOfBirth: '',
        avatar: '',
        nickname: '',
        aboutMe: ''
    });
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter(); 

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setSuccess('');
    
        try {
            const response = await apiRequest('/register', 'POST', formData);
            console.log('Registered:', response);
            setSuccess('Registration successful!');
    
            
            setTimeout(() => {
                router.push('/login'); 
            }, 1000); 
        } catch (error) {
            setError(error.message);
        } finally {
            setLoading(false);
        }
    };      

    const handleBack = () => {
        router.back(); 
    };

    return (
        <div className="container">
            <h1>Register</h1>
            {success && <p className="success">{success}</p>}
            {error && <p className="error">{error}</p>}
            <form className="registrationForm" onSubmit={handleSubmit}>
                <input type="text" name="firstName" placeholder="First Name" onChange={handleChange} required />
                <input type="text" name="lastName" placeholder="Last Name" onChange={handleChange} required />
                <input type="email" name="email" placeholder="Email" onChange={handleChange} required />
                <input type="password" name="password" placeholder="Password" onChange={handleChange} required />
                <input type="date" name="dateOfBirth" onChange={handleChange} required />
                <input type="text" name="avatar" placeholder="Avatar URL (Optional)" onChange={handleChange} />
                <input type="text" name="nickname" placeholder="Nickname (Optional)" onChange={handleChange} />
                <textarea name="aboutMe" placeholder="About Me (Optional)" onChange={handleChange}></textarea>
                <button type="submit" disabled={loading}>{loading ? 'Registering...' : 'Register'}</button>
                <button type="button" onClick={handleBack} style={{ marginTop: '10px' }}>Back</button>
            </form>
        </div>
    );
};

export default Register;
