package main

import (
	"log"
	"net/http"
	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/db"
	"social-network/backend/pkg/followers"
	"social-network/backend/pkg/notifications" // Teavituste pakett
	"social-network/backend/pkg/posts"

	"github.com/gorilla/handlers"
)

func main() {
	// Andmebaasi ühendamine ja migratsioonide käivitamine
	err := db.ConnectSQLite("database.db")
	if err != nil {
		log.Fatalf("Andmebaasiga ühendamine ebaõnnestus: %v", err)
	}
	defer db.CloseSQLite()
	db.Migrate("backend/pkg/db/migrations")

	// Route'ide seadistamine
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/profile", auth.ProfileHandler)
	http.HandleFunc("/followers", followers.FollowHandler)
	http.HandleFunc("/followers/unfollow", followers.UnfollowHandler)
	http.HandleFunc("/following", followers.GetFollowingHandler)
	http.HandleFunc("/user", auth.UsersHandler)
	http.HandleFunc("/posts", posts.CreatePost)
	http.HandleFunc("/posts/user", posts.GetPosts)
	http.HandleFunc("/posts/comments", posts.CreateComment)
	

	// Teavituste route'id
	http.HandleFunc("/notifications/unread", notifications.HandleGetUnreadNotifications) // Kasutaja lugemata teavitused
	http.HandleFunc("/notifications/read/", notifications.HandleMarkNotificationAsRead)  // Märgi teavitus loetuks

	log.Println("Route'id seadistatud edukalt.")

	// CORS haldur
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                           // Lubatud kõik päritolud
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),                      // Lubatud meetodid
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "User-Email"}), // Lubatud päised
	)

	// Serveri käivitamine
	log.Println("Server käivitati aadressil :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", corsHandler(http.DefaultServeMux)); err != nil {
		log.Fatalf("Serveri käivitamine ebaõnnestus: %v", err)
	}
}
