package main

import (
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
	"social-network/backend/pkg/followers"
	"social-network/backend/pkg/following"
	"social-network/backend/pkg/groups"
	"social-network/backend/pkg/notifications"
	"social-network/backend/pkg/posts"
    "social-network/backend/pkg/chat"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Connect to the SQLite database
	err := db.ConnectSQLite("database.db")
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.CloseSQLite()

	// Run migrations
	db.Migrate("backend/pkg/db/migrations")

	// Create a new router
	router := mux.NewRouter()

	// Route configurations
	router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	router.Handle("/profile", auth.AuthMiddleware(http.HandlerFunc(auth.ProfileHandler))).Methods("GET")

	// Followers routes
	router.Handle("/followers", auth.AuthMiddleware(http.HandlerFunc(followers.FollowHandler))).Methods("POST")
	router.Handle("/followers/unfollow", auth.AuthMiddleware(http.HandlerFunc(followers.UnfollowHandler))).Methods("POST")
	router.Handle("/followers/list", auth.AuthMiddleware(http.HandlerFunc(followers.GetFollowersHandler))).Methods("GET")
	router.HandleFunc("/followers/list/{id}", followers.GetUserFollowersHandler).Methods("GET")
	router.HandleFunc("/followers/status/{id}", followers.CheckFollowStatusHandler).Methods("GET")

	// Following routes
	router.Handle("/following/list", auth.AuthMiddleware(http.HandlerFunc(following.GetFollowingHandler))).Methods("GET")
	router.HandleFunc("/following/list/{id}", following.GetUserFollowingHandler).Methods("GET")

	// User routes
	router.HandleFunc("/users", auth.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", auth.UserProfileHandler).Methods("GET")

	// Post routes
	router.HandleFunc("/posts/user", posts.GetPosts).Methods("GET")
	router.HandleFunc("/posts", posts.CreatePost).Methods("POST")
	router.HandleFunc("/posts/comments", posts.CreateComment).Methods("POST")
	// Serve static files from the "uploads" directory
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))


	// Groups routes
	router.HandleFunc("/groups/create", groups.CreateGroup).Methods("POST")
	router.HandleFunc("/groups", groups.GetGroups).Methods("GET")

	// Notification routes
	router.Handle("/notifications/unread", auth.AuthMiddleware(http.HandlerFunc(notifications.HandleGetUnreadNotifications))).Methods("GET")
	router.Handle("/notifications/read/{id}", auth.AuthMiddleware(http.HandlerFunc(notifications.HandleMarkNotificationAsRead))).Methods("POST")

	// Chat routes
    router.HandleFunc("/ws", chat.HandleConnections)

	// CORS handler
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "User-Email"}),
	)

	// Start the server
	log.Println("Server started on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", corsHandler(router)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
