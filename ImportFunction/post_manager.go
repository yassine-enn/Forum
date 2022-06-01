package ImportFunction

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	PostID       int
	PostDate     string
	PostCategory string
	PostContent  string
	PostTitle    string
	PostLike     int
	PostAuthor   string
}

func PostDataReader(condition string) []Post {

	var postTable []Post
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	defer db.Close()
	result, err1 := db.Query(`SELECT PostID, date(Date), PostCategory, PostText, PostTitle, likeCounter, PostAuthor FROM Post WHERE ` + condition)
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
	var PostAuthor string
	for result.Next() {
		result.Scan(&PostID, &PostDate, &PostCategory, &PostText, &PostTitle, &PostLike, &PostAuthor)
		var post = Post{PostID, PostDate, PostCategory, PostText, PostTitle, PostLike, PostAuthor}
		postTable = append(postTable, post)
		fmt.Println(postTable)
	}
	result.Close()
	db.Close()
	return postTable
}

func PostTopic(postText string, postTitle string, postCategory string, author string) {
	db := BddOpener()
	statement, prepareErr := db.Prepare("INSERT INTO Post (Date, PostCategory, PostText, PostTitle, likeCounter, PostAuthor) VALUES (?,?,?,?,?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(time.Now(), postCategory, postText, postTitle, 0, author)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
}

func DeleteTopic(ID int) {
	db := BddOpener()
	statement, prepareErr := db.Prepare("DELETE FROM Post WHERE PostID = ?")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(ID)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
}
