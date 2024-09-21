package posts

import (
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

    // Extract form data
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
    rows, err := db.DB.Query("SELECT id, user_id, content, image, gif, privacy, created_at FROM posts ORDER BY created_at DESC")
    if err != nil {
        log.Printf("Error fetching posts: %v", err)
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Image, &post.GIF, &post.Privacy, &post.CreatedAt)
        if err != nil {
            log.Printf("Error scanning post: %v", err)
            continue
        }
        posts = append(posts, post)
    }

    json.NewEncoder(w).Encode(posts)
}
