package backend

import (
	"database/sql"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// --- URL Validation ---
		if len(r.URL.Path) > 2048 || strings.Contains(r.URL.Path, "\x00") || strings.Contains(r.URL.Path, "..") {
			log.Printf("Suspicious login path: %q", r.URL.Path)
			Render(w, http.StatusBadRequest)
			return
		}
		if r.URL.Path != path.Clean(r.URL.Path) {
			Render(w, http.StatusBadRequest)
			return
		}

		// --- Method Validation ---
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			Render(w, http.StatusMethodNotAllowed)
			return
		}

		// --- GET Method ---
		if r.Method == http.MethodGet {
			if err := templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": ""}); err != nil {
				log.Printf("Template render error (GET /login): %v", err)
				Render(w, http.StatusInternalServerError)
			}
			return
		}

		// --- POST Method ---
		if err := r.ParseForm(); err != nil {
			Render(w, http.StatusBadRequest)
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
		err := DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ? OR username = ?", email, email).Scan(&userID, &passwordHash)
		if err == sql.ErrNoRows {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid email or password"})
			return
		}
		if err != nil {
			log.Printf("DB error on login: %v", err)
			Render(w, http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
			templates.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid email or password"})
			return
		}

		// --- Create Session ---
		token := generateRandomToken()
		if token == "" {
			log.Printf("Token generation failed")
			Render(w, http.StatusInternalServerError)
			return
		}
		exp := time.Now().Add(24 * time.Hour)

		_, err = DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		if err != nil {
			log.Printf("Failed to remove old sessions: %v", err)
			Render(w, http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec(
			"INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)",
			token, userID, exp.Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			log.Printf("Insert session error: %v", err)
			Render(w, http.StatusInternalServerError)
			return
		}

		if err != nil {
			log.Printf("Insert session error: %v", err)
			Render(w, http.StatusInternalServerError)
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

		http.Redirect(w, r, "/post", http.StatusSeeOther)
	}
}
