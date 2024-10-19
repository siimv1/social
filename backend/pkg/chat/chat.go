package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
	"social-network/backend/pkg/followers"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Represents a WebSocket client
type Client struct {
	UserID int
	GroupID int
	Conn   *websocket.Conn
	Send   chan []byte
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  2048, 
    WriteBufferSize: 2048,  
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var userClients = make(map[int]*Client)
var groupClients = make(map[int]map[int]*Client)  

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer ws.Close()
	senderIDStr := r.URL.Query().Get("sender_id")
	recipientIDStr := r.URL.Query().Get("recipient_id")
	log.Printf("Received WebSocket connection request - Sender ID: %s, Recipient ID: %s", senderIDStr, recipientIDStr)
	if senderIDStr == "" || recipientIDStr == "" {
		log.Println("Invalid Sender or Recipient ID: IDs are missing")
		return
	}
	senderID, err := strconv.Atoi(senderIDStr)
	if err != nil {
		log.Printf("Invalid Sender ID: %v", err)
		return
	}
	recipientID, err := strconv.Atoi(recipientIDStr)
	if err != nil {
		log.Printf("Invalid Recipient ID: %v", err)
		return
	}
	log.Printf("Checking mutual follow status between user %d and user %d", senderID, recipientID)
	// Check if sender and recipient are mutually following each other
	isMutualFollow, err := followers.CheckMutualFollowStatus(senderID, recipientID)
	if err != nil {
		log.Printf("Error checking mutual follow status between user %d and user %d: %v", senderID, recipientID, err)
		ws.WriteMessage(websocket.TextMessage, []byte("Error occurred while checking follow status."))
		return
	}
	if !isMutualFollow {
		log.Printf("Private messaging is not allowed between user %d and user %d", senderID, recipientID)
		ws.WriteMessage(websocket.TextMessage, []byte("Messaging not allowed: you must be mutually following each other."))
		return
	}
	log.Printf("Users %d and %d are mutually following each other. Proceeding with chat.", senderID, recipientID)
	client := &Client{UserID: senderID, Conn: ws, Send: make(chan []byte, 256)}
	userClients[senderID] = client
	go client.WritePump()
	log.Printf("Retrieving chat history between user %d and user %d", senderID, recipientID)
	chatHistory, err := retrieveChatHistory(senderID, recipientID)
	if err == nil {
		for _, msg := range chatHistory {
			client.Send <- msg
		}
	}
	log.Printf("Listening for incoming messages between user %d and user %d", senderID, recipientID)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket closure: %v", err)
			}
			delete(userClients, senderID)
			break
		}
		handleIncomingMessage(client, message)
	}
}

// Send message only to the intended recipient
func handleIncomingMessage(sender *Client, message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}
	recipientID, ok := msg["recipient_id"].(float64)
	if !ok {
		return
	}
	content, ok := msg["content"].(string)
	if !ok {
		return
	}
	err := saveMessage(sender.UserID, int(recipientID), content)
	if err != nil {
		return
	}
	log.Printf("Message from User %d to User %d: %s", sender.UserID, int(recipientID), content)
	recipientClient, exists := userClients[int(recipientID)]
	if exists {
		recipientClient.Send <- message
	} else {
	}
}
func retrieveChatHistory(userID int, recipientID int) ([][]byte, error) {
	// Query to select messages between the user and the recipient in both directions
	rows, err := db.DB.Query(`
        SELECT sender_id, recipient_id, content, created_at
        FROM Messages
        WHERE (sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)
        ORDER BY created_at ASC
    `, userID, recipientID, recipientID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chatHistory [][]byte
	for rows.Next() {
		var senderID, recipientID int
		var content, createdAt string
		if err := rows.Scan(&senderID, &recipientID, &content, &createdAt); err != nil {
			continue
		}
		// Create a map for each message and marshal to JSON
		msg := map[string]interface{}{
			"sender_id":    senderID,
			"recipient_id": recipientID,
			"content":      content,
			"created_at":   createdAt,
		}
		messageJSON, _ := json.Marshal(msg)
		chatHistory = append(chatHistory, messageJSON)
	}
	return chatHistory, nil
}

// Save message to the database
func saveMessage(senderID int, recipientID int, content string) error {
	query := `INSERT INTO Messages (sender_id, recipient_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := db.DB.Exec(query, senderID, recipientID, content, time.Now().Format(time.RFC3339))
	return err
}
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	close(c.Send)
	delete(userClients, c.UserID)
}
func HandleGroupConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        return
    }

    groupIDStr := r.URL.Query().Get("group_id")
    userIDStr := r.URL.Query().Get("user_id")

    if groupIDStr == "" || userIDStr == "" {
        log.Println("Missing group ID or user ID in WebSocket request")
        return
    }

    groupID, _ := strconv.Atoi(groupIDStr)
    userID, _ := strconv.Atoi(userIDStr)

    log.Printf("User %d connected to group %d", userID, groupID)

    client := &Client{
        UserID:  userID,
        GroupID: groupID,
        Conn:    ws,
        Send:    make(chan []byte),
    }

    if groupClients[groupID] == nil {
        groupClients[groupID] = make(map[int]*Client)
    }
    groupClients[groupID][userID] = client

    // Only defer cleanup when connection is closed, not right after upgrade
    go handleGroupMessages(client, groupID)

    log.Printf("Successfully upgraded connection for user %d in group %d", userID, groupID)
}



func handleGroupMessages(client *Client, groupID int) {
    defer func() {
        log.Printf("Closing connection for user %d in group %d", client.UserID, groupID)
        client.Conn.Close()
        delete(groupClients[groupID], client.UserID)
    }()

    for {
        _, msg, err := client.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error reading message: %v", err)
            } else {
                log.Printf("WebSocket closed normally for user %d in group %d", client.UserID, groupID)
            }
            break // Exit the loop on error or closed connection
        }

        // Log incoming message
        log.Printf("Received message from user %d in group %d: %s", client.UserID, groupID, string(msg))

        saveGroupMessage(client.GroupID, client.UserID, string(msg))

        // Broadcast the message to all users in the group
        for _, c := range groupClients[groupID] {
            select {
            case c.Send <- msg:
                log.Printf("Sent message to user %d in group %d", c.UserID, groupID)
            default:
                log.Printf("Client %d is not available to receive messages", c.UserID)
            }
        }
    }
}


// Save a group message to the database
func saveGroupMessage(groupID, userID int, content string) {
	stmt, err := db.DB.Prepare("INSERT INTO group_messages (group_id, user_id, content) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return
	}
	_, err = stmt.Exec(groupID, userID, content)
	if err != nil {
		log.Printf("Failed to save group message: %v", err)
	}
}