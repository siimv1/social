import { useEffect, useState } from 'react';
import axios from 'axios';

const YourTimeline = () => {
    const [posts, setPosts] = useState([]);

    useEffect(() => {
        axios.get('/api/posts')
            .then(response => {
                console.log(response.data);  // Add this line to log the API response
                setPosts(response.data);
            })
            .catch(error => {
                console.error('Error fetching posts:', error);
            });
    }, []);

    return (
        <div>
            <h2>Your Timeline</h2>
            <div>
                {posts.map(post => (
                    <div key={post.id}>
                        <p>{post.content}</p>
                        {post.image && <img src={post.image} alt="post image" />}
                        {post.gif && <img src={post.gif} alt="post gif" />}
                        <p>Posted at: {new Date(post.created_at).toLocaleString()}</p>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default YourTimeline;
