import React, { useState, useEffect, useRef } from 'react';

const Chat = ({ senderId, recipientId }) => {
  const [messages, setMessages] = useState([]); // Store chat messages
  const [input, setInput] = useState(''); // Manage message input
  const ws = useRef(null); // Store WebSocket instance

useEffect(() => {
  if (senderId && recipientId) {
    const wsUrl = `ws://localhost:8080/ws?sender_id=${senderId}&recipient_id=${recipientId}`;
    console.log("WebSocket URL:", wsUrl);

    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      console.log('WebSocket connection established.');
    };

    ws.current.onmessage = (event) => {
      try {
        const parsedData = JSON.parse(event.data); // Parse the message if JSON
        console.log("Received message:", parsedData);
        setMessages((prevMessages) => [...prevMessages, parsedData]);
      } catch (e) {
        console.error("Failed to parse WebSocket message:", event.data);
      }
    };

    ws.current.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.current.onclose = () => {
      console.log("WebSocket connection closed.");
    };
  }

  return () => {
    if (ws.current) {
      console.log("Closing WebSocket connection...");
      ws.current.close(); // Ensure WebSocket is closed gracefully
    }
  };
}, [senderId, recipientId]);


  const sendMessage = () => {
    if (input.trim() && ws.current && ws.current.readyState === WebSocket.OPEN) {
      const message = {
        sender_id: senderId,
        recipient_id: Number(recipientId),
        content: input,
      };
      console.log("Sending message:", message);

      // Immediately update the UI to show the message
      setMessages((prevMessages) => [...prevMessages, message]);

      // Send the message to the WebSocket server
      ws.current.send(JSON.stringify(message));
      setInput('');
    }
  };

  return (
    <div style={styles.chatContainer}>
      <h2 style={styles.chatHeader}>Chat with User {recipientId}</h2>
      <div style={styles.chatMessages}>
        {/* Display messages */}
        {messages.map((msg, index) => (
          <div key={index} style={styles.message}>
            {msg.sender_id === senderId
              ? `You: ${msg.content}` // Show "You" for the sender
              : `User ${msg.sender_id}: ${msg.content}`} {/* Show sender's ID dynamically */}
          </div>
        ))}
      </div>
      <div style={styles.chatInputContainer}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message..."
          style={styles.chatInput}
        />
        <button onClick={sendMessage} style={styles.sendButton}>Send</button>
      </div>
    </div>
  );
};

// Define styles for the chat component
const styles = {
  chatContainer: {
    display: 'flex',
    flexDirection: 'column',
    width: '300px',
    border: '1px solid #ddd',
    borderRadius: '8px',
    backgroundColor: '#f9f9f9',
    height: '400px',
    overflow: 'hidden',
  },
  chatHeader: {
    backgroundColor: '#4267B2',
    color: 'white',
    padding: '8px',
    textAlign: 'center',
  },
  chatMessages: {
    flex: 1,
    padding: '10px',
    border: '1px solid #ddd',
    overflowY: 'scroll',
    marginBottom: '10px',
  },
  message: {
    padding: '5px',
    borderBottom: '1px solid #eee',
    wordWrap: 'break-word',
  },
  chatInputContainer: {
    display: 'flex',
    padding: '10px',
    borderTop: '1px solid #ddd',
    backgroundColor: '#fff',
  },
  chatInput: {
    flex: 1,
    padding: '10px',
    borderRadius: '4px',
    border: '1px solid #ddd',
    marginRight: '10px',
  },
  sendButton: {
    padding: '10px 15px',
    backgroundColor: '#004080',
    border: 'none',
    color: 'white',
    cursor: 'pointer',
    borderRadius: '4px',
  },
};

export default Chat;
