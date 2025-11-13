package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/backend"
)

func main() {
	DB, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()
	backend.LoadTemplates("templates/*.html")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/post", backend.HandlePost(DB))
	http.HandleFunc("/", backend.Handler(DB))

	// http.HandleFunc("/static", backend.HandlerStatic)
	http.HandleFunc("/signup", backend.SignupHandler(DB))
	http.HandleFunc("/login", backend.LoginHandler(DB))
	http.HandleFunc("/logout", backend.LogoutHandler(DB))
	http.HandleFunc("/comment", backend.HandleAddComment(DB))
	

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
