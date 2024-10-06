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

// Client represents a connected client in the chat.
type Client struct {
    UserID int
    Conn   *websocket.Conn
    Send   chan []byte
}

// WebSocket upgrader to handle HTTP -> WebSocket upgrade.
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// Track connected clients and broadcast channel.
var clients = make(map[*Client]bool)
var broadcast = make(chan []byte)

// HandleConnections upgrades the initial HTTP request to a WebSocket connection.
func HandleConnections(w http.ResponseWriter, r *http.Request) {
    log.Println("Received a connection request...")

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        return
    }
    defer ws.Close()

    senderIDStr := r.URL.Query().Get("sender_id")
    if senderIDStr == "" {
        log.Println("Invalid Sender ID: sender_id is missing")
        return
    }

    senderID, err := strconv.Atoi(senderIDStr)
    if err != nil {
        log.Printf("Invalid Sender ID: %v", err)
        return
    }

    client := &Client{UserID: senderID, Conn: ws, Send: make(chan []byte, 256)}
    clients[client] = true
    log.Printf("User %d connected successfully.", senderID)

    go client.WritePump()

    // Retrieve and send chat history to the connected user
    chatHistory, err := retrieveChatHistory(senderID)
    if err == nil {
        for _, msg := range chatHistory {
            client.Send <- msg
        }
    }

    // Read messages from the client
    for {
        _, message, err := ws.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Error reading message: %v", err)
            }
            log.Printf("User %d disconnected.", senderID)
            delete(clients, client)
            break
        }
        log.Printf("Received message from User %d: %s", senderID, message)
        handleIncomingMessage(client, message)
    }
}

func retrieveChatHistory(userID int) ([][]byte, error) {
    rows, err := db.DB.Query(`
        SELECT sender_id, recipient_id, content, created_at
        FROM Messages
        WHERE sender_id = ? OR recipient_id = ?
        ORDER BY created_at ASC
    `, userID, userID)

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


// Handle incoming messages and route them to the correct recipient.
func handleIncomingMessage(sender *Client, message []byte) {
    var msg map[string]interface{}
    if err := json.Unmarshal(message, &msg); err != nil {
        log.Printf("Error parsing message: %v", err)
        return
    }

    recipientID, ok := msg["recipient_id"].(float64)
    if !ok {
        log.Println("Invalid recipient ID in message.")
        return
    }

    content, ok := msg["content"].(string)
    if !ok {
        log.Println("Invalid content in message.")
        return
    }

    // Save the message to the database
    err := saveMessage(sender.UserID, int(recipientID), content)
    if err != nil {
        log.Printf("Failed to save message: %v", err)
        return
    }

    log.Printf("Message from User %d to User %d: %s", sender.UserID, int(recipientID), content)

    // Check if the recipient is a connected client
    for client := range clients {
        if client.UserID == int(recipientID) {
            client.Send <- message
            log.Printf("Message successfully sent from User %d to User %d.", sender.UserID, int(recipientID))
            return
        }
    }

    log.Printf("Recipient User %d not connected.", int(recipientID))
}

// Save the message to the database.
func saveMessage(senderID int, recipientID int, content string) error {
    query := `INSERT INTO Messages (sender_id, recipient_id, content, created_at) VALUES (?, ?, ?, ?)`
    _, err := db.DB.Exec(query, senderID, recipientID, content, time.Now().Format(time.RFC3339))
    return err
}
func HandleMessages() {
    for msg := range broadcast {
        for client := range clients {
            select {
            case client.Send <- msg:
            default:
                close(client.Send)
                delete(clients, client)
            }
        }
    }
}
func (c *Client) WritePump() {
    defer func() {
        log.Println("Client WritePump closing WebSocket...")
        c.Conn.Close()
    }()

    for msg := range c.Send {
        err := c.Conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            log.Printf("Error sending message: %v", err)
            break
        }
    }
    close(c.Send)
    delete(clients, c)
    log.Println("WritePump finished for the client.")
}