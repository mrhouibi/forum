package backend

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type DataComment struct {
	Username  string
	Content   string
	CreatedAt string
}
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

func tableExists(DB *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := DB.QueryRow(query, tableName)
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

func WriteCategories(DB *sql.DB) {
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

func checkuser(DB *sql.DB, userid int64) bool {
	var token string
	err := DB.QueryRow(`SELECT token FROM sessions WHERE user_id = ?`, userid).Scan(&token)
	if err == sql.ErrNoRows {
		return true
	}
	_, err = DB.Exec(`DELETE FROM sessions WHERE user_id = ?`, userid)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return false
}

func InsertCategoriId(DB *sql.DB, post_id int64, categories []string) {
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

func GetPost(DB *sql.DB, category, username string, UserId int64) []Datapost {
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
	WHERE post_categories.category_id=?`, Categorie_Id)
	}

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost
		if err := row.Scan(&post.Title, &post.Content, &post.IdPost); err != nil {
			fmt.Println(err)
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

func GetPostById(DB *sql.DB, PostId int) []Datapost {
	posts := []Datapost{}
	row, err := DB.Query(`SELECT title,content,id FROM posts WHERE id =?`, PostId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost
		if err := row.Scan(&post.Title, &post.Content, &post.IdPost); err != nil {
			fmt.Println(err)
			return nil
		}
		posts = append(posts, post)
	}
	return posts
}

func GetComment(DB *sql.DB, PostId int) []DataComment {
	Comments := []DataComment{}
	rows, err := DB.Query(`
		SELECT u.username, c.comment, c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at DESC`, PostId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var DataComments DataComment
		rows.Scan(&DataComments.Username, &DataComments.Content, &DataComments.CreatedAt)
		Comments = append(Comments, DataComments) // <- هنا أضفت append
	}
	return Comments
}

func Render(w http.ResponseWriter, status int) {
	tmp, err := template.ParseFiles("templates/errorpage.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	message := ""
	switch status {
	case 400:
		message = "Bad Request."
	case 403:
		message = "Access denied: you don’t have permission to view this resource."
	case 404:
		message = "Not Found."
	case 405:
		message = "Method Not Allowed."
	default:
		message = "Internal Server Error"
	}

	mes := Message_Error{
		Status:  status,
		Message: message,
	}

	w.WriteHeader(status)

	if err := tmp.Execute(w, mes); err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
