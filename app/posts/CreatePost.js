import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './file.css';

const CreatePost = ({ onPostCreated, userId, groupId }) => {  // Add groupId as a prop
  const [content, setContent] = useState('');
  const [privacy, setPrivacy] = useState('public');
  const [image, setImage] = useState(null);
  const [gif, setGif] = useState(null);
  const [previewImage, setPreviewImage] = useState(null);
  const [previewGif, setPreviewGif] = useState(null);

  useEffect(() => {
    console.log("CreatePost component received userId:", userId);
  }, [userId]);

  useEffect(() => {
    if (image) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreviewImage(reader.result);
      };
      reader.readAsDataURL(image);
    } else {
      setPreviewImage(null);
    }
  }, [image]);

  useEffect(() => {
    if (gif) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setPreviewGif(reader.result);
      };
      reader.readAsDataURL(gif);
    } else {
      setPreviewGif(null);
    }
  }, [gif]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formData = new FormData();

    formData.append('user_id', userId);  // The current user ID
    formData.append('content', content);
    formData.append('privacy', privacy);
    if (image) formData.append('image', image);
    if (gif) formData.append('gif', gif);

    // Append group_id only if it exists
    if (groupId) {
      formData.append('group_id', groupId);  // Add the group ID if passed
    }

    try {
      const response = await axios.post('http://localhost:8080/posts', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      if (onPostCreated) {
        onPostCreated(response.data);
      }

      setContent('');
      setPrivacy('public');
      setImage(null);
      setGif(null);
      setPreviewImage(null);
      setPreviewGif(null);
    } catch (error) {
      console.error('Error creating post:', error);
      alert('There was an error while posting. Please check the console.');
    }
  };

  return (
    <div className="create-post-container">
      <form className="create-post-form" onSubmit={handleSubmit}>
        <textarea
          className="post-textarea"
          placeholder="What's on your mind?"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          required
        />

        <div className="upload-section">
          <label className="upload-label image-upload">
            <span className="upload-icon">📷</span>
            <span className="upload-text">Upload Image</span>
            <input
              type="file"
              accept="image/*"
              onChange={(e) => setImage(e.target.files[0])}
              hidden
            />
          </label>

          <label className="upload-label gif-upload">
            <span className="upload-icon">🎞️</span>
            <span className="upload-text">Upload GIF</span>
            <input
              type="file"
              accept="image/gif"
              onChange={(e) => setGif(e.target.files[0])}
              hidden
            />
          </label>
        </div>

        <div className="preview-section">
          {previewImage && (
            <div className="preview-item">
              <img src={previewImage} alt="Preview" className="preview-image" />
              <button type="button" className="remove-button" onClick={() => setImage(null)}>
                &times;
              </button>
            </div>
          )}
          {previewGif && (
            <div className="preview-item">
              <img src={previewGif} alt="GIF Preview" className="preview-image" />
              <button type="button" className="remove-button" onClick={() => setGif(null)}>
                &times;
              </button>
            </div>
          )}
        </div>

        <select
          className="privacy-select"
          value={privacy}
          onChange={(e) => setPrivacy(e.target.value)}
        >
          <option value="public">Public</option>
          <option value="private">Private</option>
          <option value="almost-private">Almost Private</option>
        </select>

        <button type="submit" className="submit-button">
          Post
        </button>
      </form>
    </div>
  );
};

export default CreatePost;
