package posts

import (
    "strconv"
    "log"
    "time"
    "encoding/json"
    "net/http"
    "social-network/backend/pkg/db"
)

type Post struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    Image     string    `json:"image"`
    GIF       string    `json:"gif"`
    Privacy   string    `json:"privacy"`
    CreatedAt time.Time `json:"created_at"`
}


// Comment represents a comment on a post
type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    Image     string    `json:"image"`
    GIF       string    `json:"gif"`
    CreatedAt time.Time `json:"created_at"`
}


func CreatePost(w http.ResponseWriter, r *http.Request) {
    var post Post

    // Parse form data
    err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
    if err != nil {
        http.Error(w, "Unable to parse form data", http.StatusBadRequest)
        return
    }

    // Extract and log form data
    userIDStr := r.FormValue("user_id")
    post.UserID, err = strconv.Atoi(userIDStr)     // Convert user_id to int
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    post.Content = r.FormValue("content")
    post.Privacy = r.FormValue("privacy")

    // Handle image and gif uploads
    post.Image = handleImageUpload(r)
    post.GIF = handleGIFUpload(r)


    // Save post to the database
    _, err = db.DB.Exec("INSERT INTO posts (user_id, content, image, gif, privacy) VALUES (?, ?, ?, ?, ?)",
        post.UserID, post.Content, post.Image, post.GIF, post.Privacy)
    if err != nil {
        log.Printf("Error saving post: %v", err)
        http.Error(w, "Error saving post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}


func CreateComment(w http.ResponseWriter, r *http.Request) {
    var comment Comment
    err := json.NewDecoder(r.Body).Decode(&comment)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    comment.Image = handleImageUpload(r)
    comment.GIF = handleGIFUpload(r)

    _, err = db.DB.Exec("INSERT INTO comments (post_id, user_id, content, image, gif) VALUES (?, ?, ?, ?, ?)",
        comment.PostID, comment.UserID, comment.Content, comment.Image, comment.GIF)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(comment)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("user_id")
    
    rows, err := db.DB.Query(`
    SELECT p.* FROM posts p
    LEFT JOIN followers f ON p.user_id = f.followed_id
    WHERE p.user_id = ?`, userID)


    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.GIF, &post.Privacy, &post.CreatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}