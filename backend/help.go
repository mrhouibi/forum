package backend

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Datapost struct {
	Title   string
	Content string
	IdPost  int
}
type Message_Error struct {
	Status  int
	Message string
}

var CategoriesId = make(map[string]int)

func tableExists(db *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := db.QueryRow(query, tableName)
	var name string
	err := row.Scan(&name)
	return err == nil
}

func InsertCategorie() {
	categories := []string{"Technology", "Science", "Education", "Engineering", "Entertainment"}
	i := 1
	for _, categorie := range categories {
		CategoriesId[categorie] = i
		i++
	}
}

func WriteCategories() {
	categories := []string{"Technology", "Science", "Education", "Engineering", "Entertainment"}
	insertcategorie := `INSERT INTO categories(categorie) VALUES (?)`

	for _, catcategorie := range categories {
		stmt, err := DB.Prepare(insertcategorie)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = stmt.Exec(catcategorie)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

func checkuser(userid int64) bool {
	var token string
	err := DB.QueryRow(`SELECT token FROM sessions WHERE user_id = ? `, userid).Scan(&token)
	if err == sql.ErrNoRows {
		return true
	}
	_, err = DB.Exec(`DELETE FROM sessions WHERE user_id = ? `, userid)
	if err != nil {
		log.Fatal(err)
	}
	return false
}

func InsertCategoriId(post_id int64, categories []string) {
	var categorie_id int
	for _, categorie := range categories {
		err := DB.QueryRow(`SELECT id FROM categories WHERE categorie = ?`, categorie).Scan(&categorie_id)
		if err != nil {
			return
		}
		_, err = DB.Exec("INSERT INTO post_categories (post_id,category_id) VALUES (?,?)", post_id, categorie_id)
		if err != nil {
			return
		}
	}
}

func GetPost(category, username string, UserId int64) []Datapost {
	posts := []Datapost{}
	Categorie_Id := CategoriesId[category]

	var row *sql.Rows
	var err error
	if category == "" {
		row, err = DB.Query(`SELECT title,content,id FROM posts`)
	} else if category == username {
		row, err = DB.Query(`SELECT title,content,id FROM posts WHERE user_id=?`, UserId)
	} else {
		row, err = DB.Query(`SELECT posts.title,posts.content,posts.id 
	FROM posts
	JOIN post_categories ON post_categories.post_id=posts.id
	WHERE post_categories.category_id=?
	`, Categorie_Id)
	}

	if err != nil {

		log.Fatal(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost
		if err := row.Scan(&post.Title, &post.Content, &post.IdPost); err != nil {
			log.Fatal(err)
			return nil
		}
		posts = append(posts, post)

	}
	if err = row.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return posts
}

func GetPostById(PostId int) []Datapost {
	posts := []Datapost{}
	row, err := DB.Query(`SELECT title,content,id FROM posts WHERE id =?`, PostId)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost
		if err := row.Scan(&post.Title, &post.Content, &post.IdPost); err != nil {
			log.Fatal(err)
			return nil
		}
		posts = append(posts, post)

	}
	if err = row.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return posts
}

func Render(w http.ResponseWriter, status int) {
	// Parse the error template file
	tmp, err := template.ParseFiles("templates/errorpage.html")
	// Set the HTTP status code in the response
	w.WriteHeader(status)
	// If there is an error loading the template, show a simple error message
	if err != nil {
		http.Error(w, "page not found", http.StatusNotFound)
		return
	}
	// Prepare the error message based on the status code
	message := ""
	switch status {
	case 400:
		message = "Bad Request."
	case 404:
		message = "Not Found."
	case 405:
		message = "Status Method Not Allowed."
	case 403:
		message="Access denied: you don’t have permission to view this resource."
	default:
		message = "Status Internal Server Error"
	}
	// Create a struct with status and message to pass to the template
	mes := Message_Error{
		Status:  status,
		Message: message,
	}
	// Execute the template and display the error page
	tmp.Execute(w, mes)
}
