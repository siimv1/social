import React, { useState, useEffect, useRef } from 'react';
import Picker from 'emoji-picker-react';  // Import the emoji picker component

const Chat = ({ senderId, recipientId }) => {
  const [messages, setMessages] = useState([]); // Store chat messages
  const [input, setInput] = useState(''); // Manage message input
  const [showPicker, setShowPicker] = useState(false); // Toggle the emoji picker
  const ws = useRef(null); // Store WebSocket instance
  
  // Establish WebSocket connection
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

  // Handle sending messages
  const sendMessage = () => {
    if (input.trim() && ws.current && ws.current.readyState === WebSocket.OPEN) {
      const message = {
        sender_id: senderId,
        recipient_id: Number(recipientId),
        content: input,  // input now includes emojis if added
      };
      console.log("Sending message:", message);

      // Immediately update the UI to show the message
      setMessages((prevMessages) => [...prevMessages, message]);

      // Send the message to the WebSocket server
      ws.current.send(JSON.stringify(message));
      setInput(''); // Clear the input after sending the message
    }
  };

  // Handle emoji click event
  const onEmojiClick = (emojiObject, event) => {
    setInput((prevInput) => prevInput + emojiObject.emoji);  // Correctly access the emoji object
    setShowPicker(false);  // Close the emoji picker after selecting an emoji
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
        {/* Toggle Emoji Picker */}
        <button onClick={() => setShowPicker((prev) => !prev)} style={styles.emojiButton}>
          😀
        </button>
        
        {/* Display Emoji Picker outside the chatContainer */}
        {showPicker && (
          <div style={styles.fixedPickerContainer}>
            <Picker onEmojiClick={onEmojiClick} />
          </div>
        )}

        {/* Message input field */}
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message..."
          style={styles.chatInput}
        />
        
        {/* Send Button */}
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
    width: '350px',
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
  emojiButton: {
    padding: '3px',
    border: 'none',
    backgroundColor: '#fff',
    cursor: 'pointer',
    marginRight: '5px',
  },
  fixedPickerContainer: {
    position: 'fixed',  // Use fixed position to place it outside the chatContainer
    bottom: '80px',
    right: '50px',
    zIndex: 1000,
  },
};

export default Chat;
