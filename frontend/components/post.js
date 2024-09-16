import { useState, useEffect } from 'react';

const PostComponent = () => {
  const [content, setContent] = useState(''); // To store new post content
  const [posts, setPosts] = useState([]);    // To store fetched posts

  // Fetch posts when the component loads
  useEffect(() => {
    fetch('/api/posts') // Assuming '/api/posts' is the correct endpoint
      .then((res) => res.json())
      .then((data) => setPosts(data))
      .catch((error) => console.error('Error fetching posts:', error));
  }, []);

  // Handle post creation
  const handleSubmit = async (e) => {
    e.preventDefault();

    const response = await fetch('/api/posts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content }),
    });

    if (response.ok) {
      const newPost = await response.json();
      setPosts([newPost, ...posts]); // Update the post list with the new post
      setContent('');                // Clear the input field
    } else {
      console.error('Error creating post');
    }
  };

  return (
    <div>
      <h1>Create a Post</h1>
      <form onSubmit={handleSubmit}>
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="What's on your mind?"
        />
        <button type="submit">Submit</button>
      </form>

      <h2>Posts</h2>
      <ul>
        {posts.map((post) => (
          <li key={post.id}>
            <p>{post.content}</p>
            <small>Posted at {new Date(post.created_at).toLocaleString()}</small>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PostComponent;
