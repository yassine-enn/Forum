package ImportFunction

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	PostID       int
	PostTitle    string
	PostContent  string
	PostDate     string
	PostLike     int
	PostCategory string
}

func PostDataReader(condition string) []Post {

	var postTable []Post
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	defer db.Close()
	result, err1 := db.Query(`SELECT PostID, Date, PostText, PostTitle, likeCounter FROM Post WHERE PostID > 0`)
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return nil
	}
	defer result.Close()
	var PostID int
	var PostTitle string
	var PostText string
	var PostDate string
	var PostLike int
	var PostCategory string
	for result.Next() {
		result.Scan(&PostID, &PostDate, &PostTitle, &PostText, &PostLike, &PostCategory)
		var post = Post{PostID, PostText, PostTitle, PostDate, PostLike, PostCategory}
		postTable = append(postTable, post)
	}
	result.Close()
	db.Close()
	return postTable
}

func PostTopic(postText string, postTitle string, postCategory string) {
	db := BddOpener()
	statement, prepareErr := db.Prepare("INSERT INTO Post (Date, PostCategory, PostText, Image, PostTitle, likeCounter) VALUES (?,?,?,?,?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(time.Now().Format("02-01-2006"), postCategory, postText, "", postTitle, 0)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
}
