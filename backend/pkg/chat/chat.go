package chat

import (
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

// Client represents a connected client in the chat.
type Client struct {
    Conn *websocket.Conn
    Send chan []byte
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

func HandleConnections(w http.ResponseWriter, r *http.Request) {
    log.Println("Received a connection request...")

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        return
    }
    defer func() {
        log.Println("Closing WebSocket connection...")
        ws.Close()
    }()

    client := &Client{Conn: ws, Send: make(chan []byte, 256)} // Buffered channel to avoid blocking
    clients[client] = true
    log.Println("Client connected successfully.")

    // Start a goroutine to handle the outgoing messages for this client
    go client.WritePump()

    // Start reading messages
    for {
        _, message, err := ws.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            log.Printf("Client disconnected: %v", err)
            delete(clients, client)
            break
        }
        log.Printf("Received message: %s", message)
        broadcast <- message
    }
}


// HandleMessages listens for incoming messages on the broadcast channel and sends them to clients.
func HandleMessages() {
    for msg := range broadcast { // Use for range to simplify channel iteration.
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

// WritePump handles sending messages to clients.
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
    log.Println("WritePump finished for the client.")
}
