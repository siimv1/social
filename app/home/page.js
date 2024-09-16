"use client";
import Link from 'next/link';
import Image from 'next/image';
import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import './home.css';

const Home = () => {
    const router = useRouter();
    const [content, setContent] = useState('');  // To store the post content
    const [posts, setPosts] = useState([         // Initially, placeholder posts
        { id: 1, name: 'John Doe', content: 'This is a sample post in your timeline. Like, comment, or share!' },
        { id: 2, name: 'Jane Smith', content: 'Another sample post with a similar style to Facebook\'s feed.' },
    ]);

    const handleLogout = async () => {
        router.push('/login'); 
    };

    const handlePost = async () => {
        if (!content.trim()) {
            alert('Post content cannot be empty!');
            return;
        }

        try {
            // Send the POST request to the backend to create a post
            const response = await fetch('http://localhost:8080/backend/pkg/api/posts', {  // Update with your API URL
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ content }),
            });

            if (response.ok) {
                const newPost = await response.json();
                // Add the new post to the top of the timeline
                setPosts([{ id: newPost.id, name: 'You', content: newPost.content }, ...posts]);
                setContent('');  // Clear the textarea
                alert('Post created successfully!');
            } else {
                console.error('Failed to create post', response.status);
                alert('Failed to create post');
            }
        } catch (error) {
            console.error('Error creating post:', error);
        }
    };

    return (
        <div className="home-container">
            <div className="home-header">
                <Link href="/profile">
                    <Image src="/profile.png" alt="profile" width={100} height={100} className="profile-pic" />
                </Link>
                <a href="/home" style={{ textDecoration: 'none', color: 'inherit' }} className="home-title-link">
                    <h1>Social Network</h1>
                </a>
                <div className="header-buttons">
                    <button className="notification-button">
                        <Image src="/notification.png" alt="Notifications" width={40} height={40} />
                    </button>
                    <button className="messenger-button">
                        <Image src="/messenger.png" alt="Messenger" width={50} height={50} />
                    </button>
                </div>
                <button className="logout-button" onClick={handleLogout}>Log Out</button>
            </div>

            <div className="home-sidebar-left">
                <ul>
                    <li><Link href="/profile" style={{ textDecoration: 'none', color: 'inherit' }}>My profile</Link></li>
                    <li>I'm following</li>
                    <li>My followers</li>
                </ul>
            </div>

            <div className="home-sidebar-right">
                <ul>
                    <li>Groups</li>
                    <li>Chats</li>
                </ul>
            </div>

            <div className="home-content">
                <div className="post-section">
                    <h2>Create a Post</h2>
                    <textarea
                        placeholder="What's on your mind?"
                        rows="3"
                        value={content}
                        onChange={(e) => setContent(e.target.value)}  // Update post content on input change
                    ></textarea>
                    <button className="post-button" onClick={handlePost}>Post</button> {/* This is the Post button */}
                </div>

                <div className="timeline-section">
                    <h2>Your Timeline</h2>
                    {posts.map((post) => (
                        <div className="post" key={post.id}>
                            <h3>{post.name}</h3>
                            <p>{post.content}</p>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Home;
