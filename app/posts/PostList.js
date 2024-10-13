import React, { useEffect, useState } from 'react';
import axios from 'axios';

const PostList = ({ userId, groupId, newPost }) => {  // Added groupId as a prop
  const [posts, setPosts] = useState([]);
  const [commentInputs, setCommentInputs] = useState({});

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

  useEffect(() => {
    const fetchPosts = async () => {
       let url = `http://localhost:8080/posts/user?user_id=${userId}`;
       if (groupId) {
          url = `http://localhost:8080/posts?group_id=${groupId}`;  // Fetch group posts if groupId is provided
       }

       try {
          const response = await axios.get(url);
          setPosts([...new Map(response.data.map((post) => [post.id, post])).values()]);
       } catch (error) {
          console.error('Error fetching posts:', error);
       }
    };
    fetchPosts();
  }, [userId, groupId]);  // Fetch posts when userId or groupId changes

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

  const handleCommentSubmit = async (postId) => {
    if (!commentInputs[postId]) return;

    const newComment = {
      post_id: postId,
      user_id: userId,  
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
      {posts.map((post) => (
        <div key={`${post.id}-${post.created_at}`} style={postBoxStyle}>
          <h3>{post.content}</h3>
          {post.image && <img src={`http://localhost:8080/${post.image}`} alt="Post Image" style={imageStyle} />}
          {post.gif && <img src={`http://localhost:8080/${post.gif}`} alt="Post GIF" style={imageStyle} />}
          <p><strong>Privacy:</strong> {post.privacy}</p>
          {post.comments &&
            post.comments.map((comment, index) => (
              <div key={index} className="comment-container" style={commentContainerStyle}>
                <span className="comment-author">
                  {comment.first_name} {comment.last_name}:
                </span>
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
            <button
              onClick={() => handleCommentSubmit(post.id)}
              style={submitButtonStyle}
            >
              Post
            </button>
          </div>
        </div>
      ))}
    </div>
  );
};

export default PostList;

