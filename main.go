package main

import (
    "log"
    "net/http"
    "social-network/backend/pkg/auth"
    "social-network/backend/pkg/api"  // Import your post handlers here
    "social-network/backend/pkg/db"
    "github.com/gorilla/handlers"
)

func main() {
    // Connect to the SQLite database
    err := db.ConnectSQLite("database.db")
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
    }
    defer db.CloseSQLite()
    db.Migrate("backend/pkg/db/migrations")

    // Set up your routes for auth
    http.HandleFunc("/register", auth.RegisterHandler)
    http.HandleFunc("/login", auth.LoginHandler)

    // Set up routes for posts
    http.HandleFunc("/api/posts", api.CreatePost) // POST to create a post
    http.HandleFunc("/api/posts", api.GetPosts)   // GET to fetch all posts

    // Create a CORS handler
    corsHandler := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}),                             // Allow all origins
        handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),        // Allow specified methods
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow specified headers
    )
	
    // Start the server with CORS enabled
    log.Println("Server started at :8080")
    if err := http.ListenAndServe("0.0.0.0:8080", corsHandler(http.DefaultServeMux)); err != nil {
        log.Fatalf("Could not start the server: %v", err)
    }
}
