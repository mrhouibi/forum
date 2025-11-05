package backend

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": ""}); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Invalid form"})
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Println(username)
		fmt.Println(email)
		fmt.Println(password)
		// validations
		if username == "" || email == "" || password == "" {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "All fields are required"})
			return
		}
		if !emailRe.MatchString(email) {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Invalid email"})
			return
		}
		if len(password) < 8 {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Password must be >= 8 chars"})
			return
		}

		var exists int
		err := DB.QueryRow("SELECT COUNT(1) FROM users WHERE email = ?", email).Scan(&exists)
		if err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Database error"})
			return
		}
		if exists > 0 {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Email already taken"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Server error"})
			return
		}

		tx, err := DB.Begin()
		if err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Server error"})
			return
		}
		defer tx.Rollback()

		res, err := tx.Exec("INSERT INTO users (username, email, password_hash, created_at) VALUES (?, ?, ?, datetime('now'))", username, email, string(hash))
		if err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Error creating user"})
			return
		}

		userID, err := res.LastInsertId()
		if err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Server error"})
			return
		}

		token := generateRandomToken()
		if token == "" {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Failed to generate session token"})
			return
		}
		exp := time.Now().Add(24 * time.Hour)

		_, err = tx.Exec("INSERT INTO sessions(token, user_id, expires_at) VALUES (?, ?, ?)", token, userID, exp.Format("2006-01-02 15:04:05"))
		if err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Error creating session"})
			return
		}

		if err := tx.Commit(); err != nil {
			templates.ExecuteTemplate(w, "signup.html", map[string]string{"Error": "Server error"})
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

func generateRandomToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hexEncode(b)
}

func hexEncode(b []byte) string {
	const hexChars = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexChars[v>>4]
		out[i*2+1] = hexChars[v&0x0f]
	}
	return string(out)
}
