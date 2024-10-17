package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"social-network/backend/pkg/db"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	userID, err := authenticateUser(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, err := Store.Get(r, "session-name")
	if err != nil {
		// Log the error but proceed to create a new session
		log.Println("Failed to get session:", err)
		// Create a new session
		session = sessions.NewSession(Store, "session-name")
		session.IsNew = true
	}

	session.Values["user_id"] = userID
	err = session.Save(r, w)
	if err != nil {
		log.Println("Failed to save session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// authenticateUser checks the provided email and password against the database
func authenticateUser(email, password string) (int, error) {
	var userID int
	var storedPasswordHash string

	// Fetch the user ID and hashed password from the database
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid credentials")
		}
		return 0, err
	}

	// Compare the hashed password with the provided password
	match := comparePasswords(storedPasswordHash, password)
	if !match {
		return 0, errors.New("invalid credentials")
	}

	return userID, nil
}

func comparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}
