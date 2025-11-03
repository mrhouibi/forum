package main

import (
	"log"
	"net/http"

	"forum/backend"
	"forum/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the forum"))
}

func main() {
	database.InitDB("forum.db")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/signup", backend.SignupHandler)

	http.HandleFunc("/login", backend.LoginHandler)
	http.HandleFunc("/logout", backend.LogoutHandler)

	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
