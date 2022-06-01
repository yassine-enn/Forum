package ImportFunction

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	PostID      int
	PostTitle   string
	PostContent string
	PostDate    string
	PostLike    int
}

func PostDataReader() []Post {
	var postTable []Post
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	result, err1 := db.Query(`SELECT PostID, Date, PostText, PostTitle, likeCounter FROM Post WHERE PostID > 0`)
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return nil
	}
	var PostID int
	var PostTitle string
	var PostText string
	var PostDate string
	var PostLike int
	for result.Next() {
		result.Scan(&PostID, &PostDate, &PostTitle, &PostText, &PostLike)
		var post = Post{PostID, PostText, PostTitle, PostDate, PostLike}
		postTable = append(postTable, post)
	}
	result.Close()
	db.Close()
	return postTable
}

func PostTopic(postText string, postTitle string, postCategory string, postImage string) {
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base", err)
		return
	}
	statement, prepareErr := db.Prepare("INSERT INTO Post (Date, PostCategory, PostText, Image, PostTitle, likeCounter) VALUES (?,?,?,?,?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	date := string(time.Now().Format("02-01-2006"))
	_, queryErr := statement.Exec(date, "", postText, "", postTitle, 0)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
}
