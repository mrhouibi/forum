package main

import (
	"log"
	"net/http"

	"forum/backend"
)

func main() {
	backend.LoadTemplates("templates/*.html")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", backend.Handler)
	http.HandleFunc("/post", backend.HandlePost)
	http.HandleFunc("/addpost", backend.HandleAddPost)
	// http.HandleFunc("/static", backend.HandlerStatic)
	http.HandleFunc("/signup", backend.SignupHandler)
	http.HandleFunc("/login", backend.LoginHandler)
	http.HandleFunc("/logout", backend.LogoutHandler)
	
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
