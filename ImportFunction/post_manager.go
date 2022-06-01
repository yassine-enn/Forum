package ImportFunction

import (
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
<<<<<<< HEAD
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	result, err1 := db.Query(`SELECT PostID, Date, PostText, PostTitle, likeCounter FROM Post WHERE PostID > 0`)
=======
	db := BddOpener()
	result, err1 := db.Query(`SELECT PostID, Date, PostText, PostTitle, likeCounter, PostCategory FROM Post WHERE ` + condition)
>>>>>>> 35b068193f0429612d9fe07d79731cbd1afbd2e3
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return nil
	}
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

<<<<<<< HEAD
func PostTopic(postText string, postTitle string, postCategory string, postImage string) {
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base", err)
		return
	}
=======
func PostTopic(postText string, postTitle string, postCategory string) {
	db := BddOpener()
>>>>>>> 35b068193f0429612d9fe07d79731cbd1afbd2e3
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
