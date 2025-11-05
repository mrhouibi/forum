package backend

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type PostPageData struct {
	Username string
	Posts    []Datapost
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	} else if r.Method != http.MethodGet {
		return
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, nil)
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		return
	} else if r.Method != http.MethodGet {
		return
	}
	userid := GetUserIDFromRequest(r)
	username := ""
	if userid != 0 {

		err := DB.QueryRow("SELECT username FROM users WHERE id = ?", userid).Scan(&username)
		if err != nil {
			fmt.Print(err)
			return
		}

	}
	post := GetPost()

	PostPageData := &PostPageData{Username: username, Posts: post}

	tmp, err := template.ParseFiles("templates/post.html")
	if err != nil {
		return
	}
	if err = tmp.Execute(w, PostPageData); err != nil {
		fmt.Println(err)
		return
	}
}

func HandleAddPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/addpost" {
		return
	}
	userId := GetUserIDFromRequest(r)

	if userId == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		tmp, err := template.ParseFiles("templates/addpost.html")
		if err != nil {
			return
		}
		tmp.Execute(w, nil)
	} else if r.Method == http.MethodPost {

		title := r.FormValue("title")
		content := r.FormValue("content")
		// category := r.FormValue("category_ids")
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
		}
		// var category []string
		category := r.Form["category_ids"]
		fmt.Println(category)

		insrtpost := `INSERT INTO posts (title,content,user_id) VALUES (?,?,?)`
		stmt, err := DB.Prepare(insrtpost)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer stmt.Close()
		res, err := stmt.Exec(title, content, userId)
		if err != nil {
			fmt.Println(err)
			return
		}
		IdPost, err := res.LastInsertId()
		if err != nil {
			fmt.Println("Error getting last insert ID:", err)
			return
		}
		InsertCategoriId(IdPost, category)

		http.Redirect(w, r, "/post", http.StatusSeeOther)
	} else {
		return
	}
}

func HandlerStatic(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		return
	} else {
		// Check if the requested file exists and is not a directory
		info, err := os.Stat(r.URL.Path[1:])

		if err != nil {
			return
		} else if info.IsDir() {
			return
		} else {
			// Serve the static file
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	}
}
