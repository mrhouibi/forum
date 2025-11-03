package backend

import (
	"database/sql"
	"net/http"
	"time"

	"forum/database"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": ""}); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid form"})
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		if email == "" || password == "" {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Email and password required"})
			return
		}

		var userID int64
		var passwordHash string
		err := database.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email).Scan(&userID, &passwordHash)
		if err == sql.ErrNoRows {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid email or password"})
			return
		}
		if err != nil {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Database error"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid email or password"})
			return
		}

		token := generateRandomToken()
		if token == "" {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Failed to create session"})
			return
		}
		exp := time.Now().Add(24 * time.Hour)

		_, err = database.DB.Exec("INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, exp.Format("2006-01-02 15:04:05"))
		if err != nil {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Failed to create session"})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  exp,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
