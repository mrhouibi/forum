package backend

import (
	"database/sql"
	"net/http"
)

func LogoutHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err == nil {
			token := c.Value
			_, _ = DB.Exec("DELETE FROM sessions WHERE token = ?", token)
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
			// Secure: true,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
