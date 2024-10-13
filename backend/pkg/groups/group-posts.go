package groups

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "social-network/backend/pkg/db"
    "github.com/gorilla/mux"
)

// GroupPost represents a post in a group
type GroupPost struct {
    ID        int       `json:"id"`
    GroupID   int       `json:"group_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    Privacy   string    `json:"privacy"`
    CreatedAt time.Time `json:"created_at"`
}

func GetUserIDFromContext(ctx context.Context) (int, error) {
    userID, ok := ctx.Value("userID").(int)
    if !ok {
        // Hardcoded for testing
        return 1, nil
    }
    return userID, nil
}

// CreateGroupPostHandler handles creating a new post in a group
func CreateGroupPostHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    groupID, err := strconv.Atoi(vars["group_id"])
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    var post GroupPost
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Set the user ID, group ID, and current timestamp
    post.UserID = userID
    post.GroupID = groupID
    post.CreatedAt = time.Now()

    // Insert the new post into the database
    _, err = db.DB.Exec(
        "INSERT INTO group_posts (group_id, user_id, content, privacy, created_at) VALUES (?, ?, ?, ?, ?)",
        post.GroupID, post.UserID, post.Content, post.Privacy, post.CreatedAt,
    )
    if err != nil {
        http.Error(w, "Failed to create post", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

// GetGroupPostsHandler retrieves posts for a specific group
func GetGroupPostsHandler(w http.ResponseWriter, r *http.Request) {
    userID, err := GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    groupID, err := strconv.Atoi(vars["group_id"])
    if err != nil {
        http.Error(w, "Invalid group ID", http.StatusBadRequest)
        return
    }

    // Verify if the user is a member of the group
    var isMember bool
    err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)", groupID, userID).Scan(&isMember)
    if err != nil || !isMember {
        http.Error(w, "Forbidden: You are not a member of this group", http.StatusForbidden)
        return
    }

    // Retrieve the posts for the specified group
    rows, err := db.DB.Query("SELECT id, group_id, user_id, content, privacy, created_at FROM group_posts WHERE group_id = ?", groupID)
    if err != nil {
        http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []GroupPost
    for rows.Next() {
        var post GroupPost
        if err := rows.Scan(&post.ID, &post.GroupID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt); err != nil {
            http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error reading posts", http.StatusInternalServerError)
        return
    }

    // Return the posts in JSON format
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
