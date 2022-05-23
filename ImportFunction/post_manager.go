package ImportFunction

import (
	"database/sql"
	"fmt"
)

type post struct {
	// PostID      int
	PostTitle   string
	PostContent string
	// PostDate    string
	// PostAuthor  string
}

func postDataReader() []post {
	var postTable []post
	db, err := sql.Open("sqlite3", "./ALED")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	result, err1 := db.Query("SELECT PostID, PostTitle, PostContent FROM Post WHERE PostID > 0")
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return nil
	}
	var PostID int
	var PostTitle string
	var PostContent string
	var PostDate string
	var PostAuthor string
	for result.Next() {
		result.Scan(&PostID, &PostTitle, &PostContent, &PostDate, &PostAuthor)
		fmt.Println(PostID, PostTitle, PostContent, PostDate, PostAuthor)
	}
	var post = post{PostTitle, PostContent}
	postTable = append(postTable, post)
	result.Close()
	db.Close()
	return postTable
}

func PostTopic(postText string, postTitle string) {
	db, err := sql.Open("sqlite3", "./ALED")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base", err)
		return
	}
	statement, prepareErr := db.Prepare("INSERT INTO Post (PostText, PostTitle) VALUES (?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(postText, postTitle)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
	fmt.Println("postTitle:", postTitle)
	fmt.Println("postText:", postText)
}
