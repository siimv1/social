package api

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "social-network/backend/pkg/db"
)

// Post struct
type Post struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// CreatePost handles creating a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
    var post Post
    err := json.NewDecoder(r.Body).Decode(&post)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    post.CreatedAt = time.Now()
    post.UpdatedAt = time.Now()

    query := "INSERT INTO posts (user_id, content, created_at, updated_at) VALUES (?, ?, ?, ?)"
    result, err := db.DB.Exec(query, post.UserID, post.Content, post.CreatedAt, post.UpdatedAt)
    if err != nil {
        log.Printf("Could not create post: %v", err)
        http.Error(w, "Could not create post", http.StatusInternalServerError)
        return
    }

    postID, _ := result.LastInsertId()
    post.ID = int(postID)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

// GetPosts handles fetching all posts
func GetPosts(w http.ResponseWriter, r *http.Request) {
    rows, err := db.DB.Query("SELECT id, user_id, content, created_at, updated_at FROM posts")
    if err != nil {
        log.Printf("Error fetching posts: %v", err)
        http.Error(w, "Could not fetch posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.CreatedAt, &post.UpdatedAt)
        if err != nil {
            log.Printf("Error scanning post: %v", err)
            continue
        }
        posts = append(posts, post)
    }

    json.NewEncoder(w).Encode(posts)
}

// DeletePost handles deleting a post by ID
func DeletePost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    query := "DELETE FROM posts WHERE id = ?"
    _, err := db.DB.Exec(query, id)
    if err != nil {
        log.Printf("Could not delete post: %v", err)
        http.Error(w, "Could not delete post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
