"use client";
import React, { useState } from 'react'; 
import Link from 'next/link';
import '../app/global.css'; 
import { useRouter } from 'next/navigation';

const Login = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const [loading, setLoading] = useState(false);
    const router = useRouter(); 

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
            });
    
            if (response.ok) {
                const data = await response.json();
                console.log('Login data:', data);
                
                setSuccess('Login successful!');
                setTimeout(() => {
                    router.push('/home'); 
                }, 1000); 
            } else {
                const errorData = await response.json();
                setError(errorData.message || 'Login failed. Please try again.');
                console.error('Error:', errorData);
            }
        } catch (err) {
            setError('Login failed. Please try again.');
            console.error('Fetch error:', err);
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
