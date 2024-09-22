import React, { useState, useEffect } from 'react';
import axios from 'axios';

const CreatePost = ({ onPostCreated }) => {
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [image, setImage] = useState(null);
  const [gif, setGif] = useState(null);
  const [user, setUser] = useState({ first_name: '', last_name: '' });

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const token = localStorage.getItem('token'); // Make sure the token is stored in localStorage
        const response = await axios.get('http://localhost:8080/auth', {
          headers: {
            Authorization: `Bearer ${token}`, // Add the token to the Authorization header
          },
        });
        setUser(response.data);
      } catch (error) {
        console.error('Error fetching user profile:', error);
      }
    };
  
    fetchUser();
  }, []);
  

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('content', content);
    formData.append('privacy', privacy);
    if (image) formData.append('image', image);
    if (gif) formData.append('gif', gif);

    try {
      // Sending POST request to the backend
      const response = await axios.post('http://localhost:8080/posts', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      // Call the onPostCreated function to update the post list
      if (onPostCreated) {
        onPostCreated(response.data);
      }

      // Clear the form
      setContent('');
      setPrivacy('public');
      setImage(null);
      setGif(null);

    } catch (error) {
      console.error('Error creating post:', error);
      alert('There was an error while posting. Please check the console.');
    }
  };

  return (
    <div>
      {/* Display the user's name while posting */}
      <h3>{user.first_name} {user.last_name}</h3>
      <form onSubmit={handleSubmit}>
        <textarea
          placeholder="What's on your mind?"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          required
        />
        <input type="file" onChange={(e) => setImage(e.target.files[0])} />
        <input type="file" onChange={(e) => setGif(e.target.files[0])} accept="image/gif" />
        <select value={privacy} onChange={(e) => setPrivacy(e.target.value)}>
          <option value="public">Public</option>
          <option value="private">Private</option>
          <option value="almost-private">Almost Private</option>
        </select>
        <button type="submit">Post</button>
      </form>
    </div>
  );
};

export default CreatePost;
