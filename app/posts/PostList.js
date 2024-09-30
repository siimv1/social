import React, { useEffect, useState } from 'react';
import axios from 'axios';

const PostList = ({ userId, newPost }) => {
  const [posts, setPosts] = useState([]);
  const [commentInputs, setCommentInputs] = useState({});

  const postBoxStyle = {
    position: 'relative', // Added for positioning child elements
    border: '2px solid #ddd',
    borderRadius: '10px',
    padding: '15px',
    margin: '10px 0',
    boxShadow: '2px 2px 12px rgba(0, 0, 0, 0.1)',
    backgroundColor: '#fff',
  };

  const imageStyle = {
    maxWidth: '100%',
    height: 'auto',
    marginTop: '10px',
  };

  const commentContainerStyle = {
    marginTop: '20px',
    padding: '10px',
    borderTop: '1px solid #ddd',
  };

  const commentBoxStyle = {
    width: '100%',
    borderRadius: '20px',
    border: '1px solid #ddd',
    padding: '10px',
    marginBottom: '10px',
  };

  const submitButtonStyle = {
    backgroundColor: '#007BFF',
    color: '#FFFFFF',
    border: 'none',
    borderRadius: '8px',
    padding: '5px 14px',
    cursor: 'pointer',
    fontSize: '14px',
    boxShadow: '0 4px 15px rgba(0, 123, 255, 0.3)',
    backgroundImage: 'linear-gradient(90deg, #0066ff 0%, #00ccff 100%)',
  };

  // New style for the "Created At" timestamp
  const timestampStyle = {
    position: 'absolute',
    bottom: '10px',
    right: '15px',
    fontSize: '10px',
    color: '#888',
  };

  useEffect(() => {
    const fetchPosts = async () => {
      if (!userId) {
        console.error('User ID is undefined or invalid.');
        return;
      }
      try {
        const response = await axios.get(`http://localhost:8080/posts/user?user_id=${userId}`);
        setPosts([...new Map(response.data.map((post) => [post.id, post])).values()]);
      } catch (error) {
        console.error('Error fetching posts:', error);
      }
    };
    fetchPosts();
  }, [userId]);

  useEffect(() => {
    if (newPost) {
      setPosts((prevPosts) => [newPost, ...prevPosts]);
    }
  }, [newPost]);

  const handleCommentChange = (e, postId) => {
    setCommentInputs({
      ...commentInputs,
      [postId]: e.target.value,
    });
  };

  const handleCommentSubmit = (postId) => {
    if (!commentInputs[postId]) return;

    const newComment = {
      postId,
      author: 'CurrentUser', // Replace with the logged-in user's info
      content: commentInputs[postId],
    };

    setPosts((prevPosts) =>
      prevPosts.map((post) =>
        post.id === postId
          ? { ...post, comments: post.comments ? [...post.comments, newComment] : [newComment] }
          : post
      )
    );

    setCommentInputs({ ...commentInputs, [postId]: '' });
  };

  const formatDate = (dateString) => {
    const options = { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' };
    return new Date(dateString).toLocaleDateString(undefined, options);
  };

  return (
    <div>
      {posts.map((post) => (
        <div key={`${post.id}-${post.created_at}`} style={postBoxStyle}>
          <h3>{post.content}</h3>
          {post.image && <img src={`http://localhost:8080/${post.image}`} alt="Post Image" style={imageStyle} />}
          {post.gif && <img src={`http://localhost:8080/${post.gif}`} alt="Post GIF" style={imageStyle} />}
          <p><strong>Privacy:</strong> {post.privacy}</p>

          {post.comments &&
            post.comments.map((comment, index) => (
              <div key={index} className="comment-container" style={commentContainerStyle}>
                <span className="comment-author">{comment.author}:</span>
                <span className="comment-text"> {comment.content}</span>
              </div>
            ))}

          <div className="comment-input-box" style={commentContainerStyle}>
            <input
              type="text"
              value={commentInputs[post.id] || ''}
              onChange={(e) => handleCommentChange(e, post.id)}
              placeholder="Write a comment..."
              style={commentBoxStyle}
            />
            <button onClick={() => handleCommentSubmit(post.id)} style={submitButtonStyle}>
              Send
            </button>
          </div>

          <div style={timestampStyle}>
            {formatDate(post.created_at)}
          </div>
        </div>
      ))}
    </div>
  );
};

export default PostList;
