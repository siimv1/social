package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"social-network/backend/pkg/db"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var Store *sessions.CookieStore

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	secretKey := os.Getenv("SESSION_KEY")
	fmt.Println("SESSION_KEY from .env file:", secretKey) // Debugging line
	if secretKey == "" {
		panic("SESSION_KEY must be set in .env file")
	}

	// Initialize the session store
	Store = sessions.NewCookieStore([]byte(secretKey))
	// sessions.go

	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 days
		HttpOnly: true,
		Secure:   false,                // Keep false for HTTP
		SameSite: http.SameSiteLaxMode, // Use Lax for cross-site requests over HTTP
	}
}

// GetSession returns the session for a given request
func GetSession(r *http.Request) (*sessions.Session, error) {
	return Store.Get(r, "social-network-session")
}

// SaveSession saves the session for a given request and response
func SaveSession(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return session.Save(r, w)
}

// SetSessionValue sets a key-value pair in a given session
func SetSessionValue(r *http.Request, w http.ResponseWriter, key string, value interface{}) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Values[key] = value
	return SaveSession(r, w, session)
}

// GetSessionValue retrieves a value from a given session
func GetSessionValue(r *http.Request, key string) (interface{}, error) {
	session, err := GetSession(r)
	if err != nil {
		return nil, err
	}
	value, ok := session.Values[key]
	if !ok {
		return nil, http.ErrNoCookie
	}
	return value, nil
}

// ClearSession clears all data from the session
func ClearSession(r *http.Request, w http.ResponseWriter) error {
	session, err := GetSession(r)
	if err != nil {
		return err
	}
	session.Options.MaxAge = -1 // Küpsise kehtivuse tühistamine
	return SaveSession(r, w, session)
}

// SessionInfoHandler provides session information
func SessionInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Println("Failed to get session:", err)
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	log.Println("Session Values:", session.Values)

	userID, ok := session.Values["user_id"].(int)
	if !ok {
		log.Println("No user_id in session values")
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	// Optionally, fetch more user details from the database
	var user struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}
	err = db.DB.QueryRow("SELECT id, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
