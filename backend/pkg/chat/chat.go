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
	Conn   *websocket.Conn
	Send   chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var userClients = make(map[int]*Client)

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
