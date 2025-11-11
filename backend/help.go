package backend

import (
	"database/sql"
	"fmt"
	"log"
)

type Datapost struct {
	Title   string
	Content string
	IdPost int
}

var CategoriesId  =make(map[string]int)

func tableExists(db *sql.DB, tableName string) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
	row := db.QueryRow(query, tableName)
	var name string
	err := row.Scan(&name)
	return err == nil
}
func InsertCategorie(){
	
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

// func GetPost() []Datapost {
// 	posts := []Datapost{}
// 	row, err := DB.Query(`SELECT title,content FROM posts`)
// 	if err != nil {

// 		log.Fatal(err)
// 		return nil
// 	}
// 	defer row.Close()
// 	for row.Next() {
// 		var post Datapost
// 		if err := row.Scan(&post.Title, &post.Content); err != nil {
// 			log.Fatal(err)
// 			return nil
// 		}
// 		posts = append(posts, post)

// 	}
// 	if err = row.Err(); err != nil {
// 		log.Fatal(err)
// 		return nil
// 	}
// 	return posts
// }

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

func GetPost(category,username string,UserId int64) []Datapost {
	posts := []Datapost{}
	Categorie_Id := CategoriesId[category]
	
	var row *sql.Rows
	var err error
	if category == "" {
		row, err = DB.Query(`SELECT title,content,id FROM posts`)
	}else if category==username{
		row,err=DB.Query(`SELECT title,content,id FROM posts WHERE user_id=?`,UserId)
		
	}else {
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
		if err := row.Scan(&post.Title, &post.Content,&post.IdPost); err != nil {
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
	row, err := DB.Query(`SELECT title,content,id FROM posts WHERE id =?`,PostId)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer row.Close()
	for row.Next() {
		var post Datapost
		if err := row.Scan(&post.Title, &post.Content,&post.IdPost); err != nil {
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