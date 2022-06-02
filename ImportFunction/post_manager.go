package ImportFunction

import (
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

type Comment struct {
	CommentID     int
	CommentAuthor string
	CommentLike   int
	CommentDate   string
	CommentText   string
	CommentSource int
}

func PostDataReader(condition string) []Post {
	var postTable []Post
	db := BddOpener()
	defer db.Close()
	result, err1 := db.Query(`SELECT PostID, date(Date), PostCategory, PostText, PostTitle, likeCounter, PostAuthor FROM Post WHERE ` + condition + ` ORDER BY PostID DESC`)
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
	}
	result.Close()
	db.Close()
	return postTable
}

func CommentDataReader(condition string) []Comment {
	var commentTable []Comment
	db := BddOpener()
	result, err1 := db.Query(`SELECT CommentID, CommentAuthor, CommentLike, date(CommentDate), CommentText, CommentSource FROM Comment WHERE ` + condition)
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return nil
	}
	defer result.Close()
	var CommentID int
	var CommentAuthor string
	var CommentLike int
	var CommentDate string
	var CommentText string
	var CommentSource int
	for result.Next() {
		result.Scan(&CommentID, &CommentAuthor, &CommentLike, &CommentDate, &CommentText, &CommentSource)
		var comment = Comment{CommentID, CommentAuthor, CommentLike, CommentDate, CommentText, CommentSource}
		commentTable = append(commentTable, comment)
	}
	result.Close()
	db.Close()
	fmt.Println(commentTable)
	return commentTable
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

func CommentTopic(commentText string, postID int, author string) {
	db := BddOpener()
	statement, prepareErr := db.Prepare("INSERT INTO Comment (CommentAuthor, CommentLike, CommentDate, CommentText, CommentSource) VALUES (?,?,?,?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(author, 0, time.Now(), commentText, postID)
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
