// In auth/logout.go or similar file
package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Println("Failed to get session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1 // Delete session
	err = session.Save(r, w)
	if err != nil {
		log.Println("Failed to delete session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}
