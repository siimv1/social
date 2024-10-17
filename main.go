package main

import (
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/chat"
	"social-network/backend/pkg/db"
	"social-network/backend/pkg/followers"
	"social-network/backend/pkg/following"
	"social-network/backend/pkg/groups"
	"social-network/backend/pkg/notifications"
	"social-network/backend/pkg/posts"

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
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
	router.Handle("/profile", auth.AuthMiddleware(http.HandlerFunc(auth.ProfileHandler))).Methods("GET")
	router.Handle("/profile/visibility", auth.AuthMiddleware(http.HandlerFunc(auth.UpdateProfileVisibilityHandler))).Methods("POST")
	router.HandleFunc("/session", auth.SessionInfoHandler).Methods("GET")
	router.Handle("/protected", auth.AuthMiddleware(http.HandlerFunc(auth.SomeProtectedHandler))).Methods("GET")

	// Followers routes
	router.Handle("/followers", auth.AuthMiddleware(http.HandlerFunc(followers.FollowHandler))).Methods("POST")
	router.Handle("/followers/unfollow", auth.AuthMiddleware(http.HandlerFunc(followers.UnfollowHandler))).Methods("POST")
	router.Handle("/followers/list", auth.AuthMiddleware(http.HandlerFunc(followers.GetFollowersHandler))).Methods("GET")
	router.HandleFunc("/followers/list/{id}", followers.GetUserFollowersHandler).Methods("GET")
	router.HandleFunc("/followers/status/{id}", followers.CheckFollowStatusHandler).Methods("GET")

	router.Handle("/followers/requests", auth.AuthMiddleware(http.HandlerFunc(followers.GetPendingFollowRequestsHandler))).Methods("GET")
	router.Handle("/followers/accept", auth.AuthMiddleware(http.HandlerFunc(followers.AcceptFollowRequestHandler))).Methods("POST")
	router.Handle("/followers/reject", auth.AuthMiddleware(http.HandlerFunc(followers.RejectFollowRequestHandler))).Methods("POST")

	// Following routes
	router.Handle("/following/list", auth.AuthMiddleware(http.HandlerFunc(following.GetFollowingHandler))).Methods("GET")
	router.HandleFunc("/following/list/{id}", following.GetUserFollowingHandler).Methods("GET")

	// User routes
	router.HandleFunc("/users", auth.GetAllUsersHandler).Methods("GET")
	router.HandleFunc("/users/{id}", auth.UserProfileHandler).Methods("GET")

	// Post routes
	router.HandleFunc("/posts/user", posts.GetPosts).Methods("GET")
	    router.HandleFunc("/posts", posts.GetPosts).Methods("GET")  

	router.HandleFunc("/posts", posts.CreatePost).Methods("POST")
	router.HandleFunc("/posts/comments", posts.CreateComment).Methods("POST")

	// Serve static files from the "uploads" directory
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	// Groups routes
	router.HandleFunc("/groups/create", groups.CreateGroup).Methods("POST")
	router.HandleFunc("/groups", groups.GetGroups).Methods("GET")
	router.HandleFunc("/groups/{id}", groups.GetGroupByID).Methods("GET")
	router.HandleFunc("/groups/{id}/join-request", groups.RequestJoinGroup).Methods("POST")
	router.HandleFunc("/groups/{id}/join-requests", groups.GetJoinRequests).Methods("GET")
	router.HandleFunc("/groups/{id}/join-status", groups.JoinGroup).Methods("GET")
	router.HandleFunc("/groups/{id}/join-requests/{userId}/accept", groups.AcceptJoinRequest).Methods("POST")
	router.HandleFunc("/groups/{id}/join-requests/{userId}/deny", groups.DenyJoinRequest).Methods("POST")
	router.HandleFunc("/groups/{id}/join-status", groups.JoinStatus).Methods("GET")
	router.HandleFunc("/events", groups.CreateEvent).Methods("POST")
	router.HandleFunc("/eventinvites", groups.GetEventInvites).Methods("GET")
	router.HandleFunc("/eventinvites/{inviteId}/accept", groups.AcceptEventInvite).Methods("POST")
	router.HandleFunc("/eventinvites/{inviteId}/decline", groups.DeclineEventInvite).Methods("POST")
	router.HandleFunc("/events/{eventId}/invite", groups.InviteUserToEvent).Methods("POST")

	// Notification routes
	router.Handle("/notifications/unread", auth.AuthMiddleware(http.HandlerFunc(notifications.HandleGetUnreadNotifications))).Methods("GET")
	router.Handle("/notifications/read/{id}", auth.AuthMiddleware(http.HandlerFunc(notifications.HandleMarkNotificationAsRead))).Methods("POST")

	// Chat routes
	router.HandleFunc("/ws", chat.HandleConnections)

	// CORS handler
	// CORS handler
	// main.go

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Ensure this matches your frontend origin
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "User-Email"}),
		handlers.AllowCredentials(), // Important for sessions and cookies
	)

	// Start the server
	log.Println("Server started on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", corsHandler(router)); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
