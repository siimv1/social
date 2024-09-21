package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/backend/pkg/db"
)

// UserList defineerib kasutaja struktuuri
type UserList struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"` // Muudetud nickname'iks
	IsOnline  bool   `json:"is_online"`
}

// UsersHandler on HTTP käitleja, mis tagastab kasutajate loendi
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Toome kõik kasutajad andmebaasist, sealhulgas ees- ja perekonnanimi
	rows, err := db.DB.Query("SELECT id, email, first_name, last_name, nickname, is_online FROM users") // Kasutame õigeid veerge
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []UserList
	for rows.Next() {
		var user UserList
		// Lisatud ees- ja perekonnanime lugemine
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Nickname, &user.IsOnline); err != nil {
			log.Printf("Error scanning user: %v", err)
			http.Error(w, "Failed to read user data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Tagastame kasutajate nimekirja JSON formaadis
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}
