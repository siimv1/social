package groups

import (
    "encoding/json"
    "net/http"
    "social-network/backend/pkg/db"
    "strconv"
    "time"

    "github.com/gorilla/mux"
)

// GroupComment represents a comment on a group post
type GroupComment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}



// CreateGroupCommentHandler handles creating a new comment on a group post
func CreateGroupCommentHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["post_id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    var comment GroupComment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Set the user ID, post ID, and timestamp
    comment.UserID = userID
    comment.PostID = postID
    comment.CreatedAt = time.Now()

    // Insert the new comment into the database
    _, err = db.DB.Exec(
        "INSERT INTO group_comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
        comment.PostID, comment.UserID, comment.Content, comment.CreatedAt,
    )
    if err != nil {
        http.Error(w, "Failed to create comment", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// GetGroupCommentsHandler retrieves comments for a specific group post
func GetGroupCommentsHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    postID, err := strconv.Atoi(vars["post_id"])
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    // Check if the user is authorized to view the post (i.e., they are a group member)
    var groupID int
    err = db.DB.QueryRow("SELECT group_id FROM group_posts WHERE id = ?", postID).Scan(&groupID)
    if err != nil {
        http.Error(w, "Post not found", http.StatusNotFound)
        return
    }

    // Verify if the user is a member of the group
    var isMember bool
    err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)", groupID, userID).Scan(&isMember)
    if err != nil || !isMember {
        http.Error(w, "Forbidden: You are not a member of this group", http.StatusForbidden)
        return
    }

    // Retrieve the comments for the specified post
    rows, err := db.DB.Query("SELECT id, post_id, user_id, content, created_at FROM group_comments WHERE post_id = ?", postID)
    if err != nil {
        http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var comments []GroupComment
    for rows.Next() {
        var comment GroupComment
        if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
            http.Error(w, "Failed to parse comments", http.StatusInternalServerError)
            return
        }
        comments = append(comments, comment)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error reading comments", http.StatusInternalServerError)
        return
    }

    // Return the comments in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}
