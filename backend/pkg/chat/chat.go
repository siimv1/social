package chat

import (
    "log"
    "net/http"
    "encoding/json"
    "strconv"
    "time"
    "github.com/gorilla/websocket"
    "social-network/backend/pkg/db"
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

// Handle WebSocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        return
    }
    defer ws.Close()

    senderIDStr := r.URL.Query().Get("sender_id")
    recipientIDStr := r.URL.Query().Get("recipient_id") 

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

    client := &Client{UserID: senderID, Conn: ws, Send: make(chan []byte, 256)}
    userClients[senderID] = client

    go client.WritePump()

    // Retrieve and send chat history between the sender and recipient
    chatHistory, err := retrieveChatHistory(senderID, recipientID)
    if err == nil {
        for _, msg := range chatHistory {
            client.Send <- msg
        }
    }

    // Read and handle incoming messages
    for {
        _, message, err := ws.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error reading message: %v", err)
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
