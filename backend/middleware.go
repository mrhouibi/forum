package backend

import (
	"net/http"
	"time"

	"forum/database"
)

func GetUserIDFromRequest(r *http.Request) int64 {
	c, err := r.Cookie("session_token")
	if err != nil {
		return 0
	}
	token := c.Value

	var userID int64
	var expiresAtStr string
	err = database.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAtStr)
	if err != nil {
		return 0
	}

	// parse expiresAt
	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		_, _ = database.DB.Exec("DELETE FROM sessions WHERE token = ?", token)
		return 0
	}
	if time.Now().After(expiresAt) {
		_, _ = database.DB.Exec("DELETE FROM sessions WHERE token = ?", token)
		return 0
	}
	return userID
}

func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if GetUserIDFromRequest(r) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
