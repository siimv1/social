import React, { useState } from 'react';
import axios from 'axios';

const CreatePost = ({ onPostCreated }) => {
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [image, setImage] = useState(null);
  const [gif, setGif] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('content', content);
    formData.append('privacy', privacy);
    if (image) formData.append('image', image);
    if (gif) formData.append('gif', gif);

    try {
      const response = await axios.post('http://localhost:8080/posts', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      onPostCreated(response.data);
    } catch (error) {
      console.error('Error creating post:', error);
    }
  };

  return (
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
  );
};

export default CreatePost;
