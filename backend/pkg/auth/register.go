package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DateOfBirth  string `json:"date_of_birth"`
	Avatar       string `json:"avatar,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
	AboutMe      string `json:"about_me,omitempty"`
	Password     string `json:"password"`
	IsOnline     bool   `json:"is_online"`
	IsPublic     bool   `json:"is_public"`
	IsPrivate    bool   `json:"is_private"`
	FollowStatus string `json:"follow_status"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.DB.Exec(query, user.Email, hashedPassword, user.FirstName, user.LastName, user.DateOfBirth, user.Avatar, user.Nickname, user.AboutMe)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
