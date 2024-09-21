import React, { useEffect, useState } from 'react';
import axios from 'axios';

const PostList = ({ newPost }) => {
  const [posts, setPosts] = useState([]);

  // Fetch posts from the backend
  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const response = await axios.get('http://localhost:8080/posts/user?user_id=1');
        setPosts(response.data);
      } catch (error) {
        console.error('Error fetching posts:', error);
      }
    };
    fetchPosts();
  }, []);

  // Add the new post to the list if it exists
  useEffect(() => {
    if (newPost) {
      setPosts(prevPosts => [newPost, ...prevPosts]);
    }
  }, [newPost]);

  return (
    <div>
      {posts.map((post) => (
        <div key={post.id} 
          style={{
            border: '1px solid #ccc',
            borderRadius: '8px',
            padding: '16px',
            marginBottom: '16px',
            backgroundColor: '#f9f9f9',
            boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
          }}
        >
          {/* Display the first_name and last_name */}
          <h3>{post.first_name} {post.last_name}</h3>
          <p>{post.content}</p>
          {post.image && <img src={`http://localhost:8080/${post.image}`} alt="Post Image" />}
          {post.gif && <img src={`http://localhost:8080/${post.gif}`} alt="Post GIF" />}
          <p><strong>Privacy:</strong> {post.privacy}</p>
        </div>
      ))}
    </div>
  );
};

export default PostList;
