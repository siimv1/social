import React, { useEffect, useState } from 'react';
import axios from 'axios';

const PostList = ({ newPost }) => {
  const [posts, setPosts] = useState([]);

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

  // Prepend new post to the posts list if a new post exists
  const allPosts = newPost ? [newPost, ...posts] : posts;

  return (
    <div>
      {allPosts.map((post) => (
        <div key={post.id}>
          <h3>{post.content}</h3>
          {post.image && <img src={`http://localhost:8080/${post.image}`} alt="Post Image" />}
          {post.gif && <img src={`http://localhost:8080/${post.gif}`} alt="Post GIF" />}
          <p><strong>Privacy:</strong> {post.privacy}</p>
        </div>
      ))}
    </div>
  );
};

export default PostList;
