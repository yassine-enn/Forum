package ImportFunction

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func BddReader(wichTable string, condition string) []Post {
	var posts []Post
	db, err := sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		fmt.Println("Error when opening the BDD:", err)
		return nil
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM " + wichTable + " WHERE " + condition)
	if err != nil {
		fmt.Println("Error when reading the BDD:", err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var post Post
		rows.Scan(&post.PostID, &post.PostTitle, &post.PostContent, &post.PostDate)
		posts = append(posts, post)
	}
	return posts
}
