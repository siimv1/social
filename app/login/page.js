"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter();

    const fetchSession = async () => {
        try {
          const response = await fetch('http://localhost:8080/session', {
            method: 'GET',
            credentials: 'include', // Include cookies in the request
          });
          if (!response.ok) {
            const errorMessage = await response.text();
            throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorMessage}`);
          }
          const data = await response.json();
          // Use the session data as needed
        } catch (error) {
          console.error('Failed to fetch session:', error);
          // Handle the error appropriately
        }
      };
      

    const handleLogin = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        try {
            const response = await fetch('http://localhost:8080/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
                credentials: 'include', // Include credentials to send cookies
            });

            if (response.ok) {
                const data = await response.json();
                console.log('Login response data:', data);

                setSuccess('Login successful!');
                setTimeout(() => {
                    router.push('/home');
                }, 1000);
            } else {
                const errorData = await response.json();
                setError(errorData.message || 'Login failed. Please try again.');
            }
        } catch (parseError) {
            setError('Login failed. Please try again.');
        } finally {
            setLoading(false);
        }
    };
 
    return (
        <div className="container">
            <div className="header">
                <h1>Welcome to the Social Network</h1>
                <p>Please login to continue.</p>
            </div>

            {error && <p className="error">{error}</p>}
            {success && <p className="success">{success}</p>}

            <div className="loginForm">
                <form onSubmit={handleLogin}>
                    <input
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Email"
                        required
                    />
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="Password"
                        required
                    />
                    <button type="submit" disabled={loading}>
                        {loading ? 'Logging in...' : 'Login'}
                    </button>
                </form>
            </div>
            <p>You don't have an account? <Link href="/register">Register</Link></p>
        </div>
    );
};

export default Login;
