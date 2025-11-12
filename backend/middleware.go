package backend

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func GetUserIDFromRequest(DB *sql.DB,r *http.Request) int64 {
	c, err := r.Cookie("session_token")
	if err != nil {
		return 0
	}
	token := c.Value

	var userID int64
	var expiresAtStr string
	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAtStr)
	if err != nil {
		return 0
	}
	// parse expiresAt
	expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		fmt.Println(err)
		_, _ = DB.Exec("DELETE FROM sessions WHERE token = ?", token)
		return 0
	}
	if time.Now().After(expiresAt) {
		_, _ = DB.Exec("DELETE FROM sessions WHERE token = ?", token)
		return 0
	}
	return userID
}

func AuthRequired(DB *sql.DB,next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if GetUserIDFromRequest(DB,r) == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
