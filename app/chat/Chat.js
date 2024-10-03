import React, { useEffect, useState, useRef } from 'react';

const Chat = () => {
  const [messages, setMessages] = useState([]); // Chat messages state
  const [input, setInput] = useState(""); // Message input state
  const wsRef = useRef(null); // Ref to store WebSocket instance

  useEffect(() => {
    // Initialize WebSocket connection when component mounts
    const socket = new WebSocket("ws://localhost:8080/ws");
    wsRef.current = socket;

    socket.onopen = () => {
      console.log("WebSocket connection established.");
    };

    socket.onmessage = (event) => {
      console.log("Received message:", event.data);
      setMessages((prevMessages) => [...prevMessages, event.data]);
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    socket.onclose = (event) => {
      if (event.code !== 1000) { // 1000 = Normal Closure
        console.warn(`WebSocket closed unexpectedly: [${event.code}] ${event.reason}`);
      } else {
        console.log("WebSocket connection closed normally.");
      }
    };

    // Cleanup WebSocket connection on component unmount
    return () => {
      if (wsRef.current) {
        console.log("Cleaning up WebSocket...");
        wsRef.current.close();
      }
    };
  }, []);

  // Function to handle sending a message
  const sendMessage = () => {
    if (input.trim() && wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log("Sending message:", input);
      wsRef.current.send(input); // Send the input value through WebSocket
      setInput(""); // Clear the input field after sending the message
    } else {
      console.warn("WebSocket not open, cannot send message.");
    }
  };

  return (
    <div style={styles.chatContainer}>
      <h2 style={styles.chatHeader}>Chat</h2>
      <div style={styles.chatMessages}>
        {/* Display messages */}
        {messages.map((msg, index) => (
          <div key={index} style={styles.message}>
            {msg}
          </div>
        ))}
      </div>

      {/* Chat input and button at the bottom */}
      <div style={styles.chatInputContainer}>
        <input
          type="text"
          value={input} // Bound to the input state
          onChange={(e) => setInput(e.target.value)} // Update state when typing
          placeholder="Type a message..."
          style={styles.chatInput} // Apply style for input
        />
        <button onClick={sendMessage} style={styles.sendButton}>Send</button> {/* Apply style for button */}
      </div>
    </div>
  );
};

// Inline CSS Styles
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
    backgroundColor: '#4CAF50',
    color: 'white',
    padding: '10px',
    textAlign: 'center',
    borderTopLeftRadius: '8px',
    borderTopRightRadius: '8px',
  },
  chatMessages: {
    flex: 1, // Takes up available space
    padding: '10px',
    border: '1px solid #ddd',
    overflowY: 'scroll',
    marginBottom: '10px', // Add some space before the input area
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
    backgroundColor: '#4CAF50',
    border: 'none',
    color: 'white',
    cursor: 'pointer',
    borderRadius: '4px',
  },
};

export default Chat;
