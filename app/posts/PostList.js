import React, { useEffect, useState } from 'react';
import axios from 'axios';

axios.defaults.withCredentials = true; // Ensure cookies are sent with requests

const PostList = ({ userId, newPost, groupId }) => { // Add groupId as a prop
  const [posts, setPosts] = useState([]);
  const [commentInputs, setCommentInputs] = useState({});
  const [viewerId, setViewerId] = useState(null); // For the logged-in user ID

  const postBoxStyle = {
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
    backgroundColor: '#1877f2',
    color: 'white',
    border: 'none',
    borderRadius: '20px',
    padding: '10px 15px',
    cursor: 'pointer',
  };

  // Fetch session data on mount to retrieve viewer ID
  useEffect(() => {
    axios.get('http://localhost:8080/session', { withCredentials: true })
      .then(response => {
        setViewerId(response.data.user_id);
      })
      .catch(error => {
        console.error('Session not active:', error);
      });
  }, []);

  useEffect(() => {
    if (userId && viewerId !== null) {
      const fetchPosts = async () => {
        try {
          const response = await axios.get('http://localhost:8080/posts/user', {
            params: {
              user_id: userId,
              viewer_id: viewerId,
              group_id: groupId, // Include groupId in the API request
            },
            withCredentials: true,
          });
  
          console.log('Posts data received:', response.data); // Log the data received
  
          if (response.data && Array.isArray(response.data)) {
            const uniquePosts = response.data.filter((v, i, a) => a.findIndex(t => (t.id === v.id)) === i);
            console.log('Unique posts:', uniquePosts); // Log to see if duplicates were removed
            setPosts(uniquePosts);
          } else {
            console.log('No posts or invalid data:', response.data);
            setPosts([]);
          }
        } catch (error) {
          console.error('Failed to fetch posts:', error);
        }
      };
  
      fetchPosts();
    }
  }, [userId, viewerId, groupId]);  // Include groupId as a dependency
  

  

  const handleCommentChange = (e, postId) => {
    setCommentInputs({
      ...commentInputs,
      [postId]: e.target.value,
    });
  };

  const handleCommentSubmit = async (postId) => {
    if (!commentInputs[postId]) return;

    const newComment = {
      post_id: postId,
      user_id: viewerId,  
      content: commentInputs[postId],
    };

    try {
      const response = await axios.post('http://localhost:8080/posts/comments', newComment);
      const addedComment = response.data;

      setPosts((prevPosts) =>
        prevPosts.map((post) =>
          post.id === postId
            ? { ...post, comments: post.comments ? [...post.comments, addedComment] : [addedComment] }
            : post
        )
      );
      setCommentInputs({ ...commentInputs, [postId]: '' }); 
    } catch (error) {
      console.error('Error adding comment:', error);
    }
  };

  return (
    <div>
        {posts.length > 0 ? (
    posts.map((post) => {
        console.log("Rendering post with ID:", post.id); // This will show you what IDs are being processed
        return (
            <div key={post.id} style={postBoxStyle}>
                <h3>{post.content}</h3>
                {post.image && <img src={`http://localhost:8080/${post.image}`} alt="Post Image" style={imageStyle} />}
                {post.gif && <img src={`http://localhost:8080/${post.gif}`} alt="Post GIF" style={imageStyle} />}
                <p><strong>Privacy:</strong> {post.privacy}</p>
                {post.comments && post.comments.map((comment) => (
                    <div key={comment.id} style={commentContainerStyle}>
                        <span>{comment.first_name} {comment.last_name}: {comment.content}</span>
                    </div>
                ))}
                <input
                    type="text"
                    value={commentInputs[post.id] || ''}
                    onChange={(e) => handleCommentChange(e, post.id)}
                    placeholder="Write a comment..."
                    style={commentBoxStyle}
                />
                <button onClick={() => handleCommentSubmit(post.id)} style={submitButtonStyle}>
                    Post Comment
                </button>
            </div>
        );
    })
) : (
    <p>No posts available.</p>
)}

    </div>
);
};

export default PostList;
