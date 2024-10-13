package posts

import (
    "database/sql"
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
    GroupID   int       `json:"group_id"`   
    Content   string    `json:"content"`
    Image     string    `json:"image"`
    GIF       string    `json:"gif"`
    Privacy   string    `json:"privacy"`
    CreatedAt time.Time `json:"created_at"`
    Comments  []Comment `json:"comments"`
}

type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    GroupID   int       `json:"group_id,omitempty"`  
    Content   string    `json:"content"`
    Image     string    `json:"image"`
    GIF       string    `json:"gif"`
    CreatedAt time.Time `json:"created_at"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
    var post Post

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form data", http.StatusBadRequest)
        return
    }

    userIDStr := r.FormValue("user_id")
    post.UserID, err = strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    groupIDStr := r.FormValue("group_id")
    if groupIDStr != "" {
        post.GroupID, err = strconv.Atoi(groupIDStr)
        if err != nil {
            http.Error(w, "Invalid group ID", http.StatusBadRequest)
            return
        }
    } else {
        post.GroupID = 0  
    }

    post.Content = r.FormValue("content")
    post.Privacy = r.FormValue("privacy")
    post.Image = handleImageUpload(r)  
    post.GIF = handleGIFUpload(r)

    query := `INSERT INTO posts (user_id, group_id, content, image, gif, privacy) VALUES (?, ?, ?, ?, ?, ?)`
    _, err = db.DB.Exec(query, post.UserID, post.GroupID, post.Content, post.Image, post.GIF, post.Privacy)
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

    groupIDStr := r.URL.Query().Get("group_id")
    if groupIDStr != "" {
        comment.GroupID, err = strconv.Atoi(groupIDStr)
        if err != nil {
            http.Error(w, "Invalid group ID", http.StatusBadRequest)
            return
        }
    }

    comment.Image = handleImageUpload(r)
    comment.GIF = handleGIFUpload(r)

    result, err := db.DB.Exec("INSERT INTO comments (post_id, user_id, group_id, content, image, gif) VALUES (?, ?, ?, ?, ?, ?)",
        comment.PostID, comment.UserID, comment.GroupID, comment.Content, comment.Image, comment.GIF)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    lastInsertID, err := result.LastInsertId()
    if err != nil {
        http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
        return
    }

    row := db.DB.QueryRow("SELECT id, post_id, user_id, group_id, content, image, gif, created_at FROM comments WHERE id = ?", lastInsertID)
    var newComment Comment
    err = row.Scan(&newComment.ID, &newComment.PostID, &newComment.UserID, &newComment.GroupID, &newComment.Content, &newComment.Image, &newComment.GIF, &newComment.CreatedAt)
    if err != nil {
        http.Error(w, "Error retrieving comment details", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newComment)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
    groupIDStr := r.URL.Query().Get("group_id")
    userIDStr := r.URL.Query().Get("user_id")  
    var rows *sql.Rows
    var err error

    // Fetch posts by group ID
    if groupIDStr != "" {
        groupID, err := strconv.Atoi(groupIDStr)
        if err != nil {
            http.Error(w, "Invalid group ID", http.StatusBadRequest)
            return
        }
        rows, err = db.DB.Query(`
            SELECT id, user_id, COALESCE(group_id, 0), content, image, gif, privacy, created_at 
            FROM posts 
            WHERE group_id = ?`, groupID)
        if err != nil {
            http.Error(w, "Error fetching group posts", http.StatusInternalServerError)
            return
        }
    // Fetch posts by user ID
    } else if userIDStr != "" {
        userID, err := strconv.Atoi(userIDStr)
        if err != nil {
            http.Error(w, "Invalid user ID", http.StatusBadRequest)
            return
        }
        rows, err = db.DB.Query(`
            SELECT id, user_id, COALESCE(group_id, 0), content, image, gif, privacy, created_at 
            FROM posts 
            WHERE (user_id = ? OR privacy = 'public')`, userID)
        if err != nil {
            http.Error(w, "Error fetching user posts", http.StatusInternalServerError)
            return
        }
    } else {
        rows, err = db.DB.Query(`
            SELECT id, user_id, COALESCE(group_id, 0), content, image, gif, privacy, created_at 
            FROM posts 
            WHERE privacy = 'public'`)
        if err != nil {
            http.Error(w, "Error fetching public posts", http.StatusInternalServerError)
            return
        }
    }

    defer rows.Close()  

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.UserID, &post.GroupID, &post.Content, &post.Image, &post.GIF, &post.Privacy, &post.CreatedAt)
        if err != nil {
            log.Printf("Error scanning post data: %v", err)
            http.Error(w, "Error scanning post data", http.StatusInternalServerError)
            return
        }

        commentRows, commentErr := db.DB.Query(`
            SELECT c.id, c.post_id, c.user_id, COALESCE(c.group_id, 0), c.content, c.image, c.gif, c.created_at, u.first_name, u.last_name
            FROM comments c
            JOIN users u ON c.user_id = u.id
            WHERE c.post_id = ?`, post.ID)
        if commentErr != nil {
            http.Error(w, "Error fetching comments", http.StatusInternalServerError)
            return
        }
        defer commentRows.Close()  

        var comments []Comment
        for commentRows.Next() {
            var comment Comment
            err := commentRows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.GroupID, &comment.Content, &comment.Image, &comment.GIF, &comment.CreatedAt, &comment.FirstName, &comment.LastName)
            if err != nil {
                log.Printf("Error scanning comment data: %v", err)
                http.Error(w, "Error scanning comment data", http.StatusInternalServerError)
                return
            }

            comments = append(comments, comment)
        }
        post.Comments = comments
        posts = append(posts, post)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
