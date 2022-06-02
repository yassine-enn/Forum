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

type Category struct {
	CategoryID   int
	CategoryName string
}

// It returns a slice of Post structs, which are the posts that match the condition and the pagination
func PostDataReader(condition string, pagin string) []Post {
	var postTable []Post
	db := BddOpener()
	defer db.Close()
	result, err1 := db.Query(`SELECT PostID, date(Date), PostCategory, PostText, PostTitle, likeCounter, PostAuthor FROM Post WHERE ` + condition + ` ORDER BY PostID DESC ` + pagin)
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
	fmt.Println(postTable)
	return postTable
}

// It returns a slice of Comment structs, which are the comments that match the condition
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

// It takes a post text, a post title, a post category and an author as parameters, and inserts them
// into the database
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

// It takes a comment text, a post ID, and an author, and inserts a new comment into the database
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

// It deletes a topic from the database
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

// It opens a connection to the database, prepares a query, executes the query and closes the
// connection
func AddCategory(categoryName string) {
	db := BddOpener()
	statement, prepareErr := db.Prepare("INSERT INTO Category (CategoryName) VALUES (?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(categoryName)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
}

// It returns a slice of Category structs, which are read from the database
func CategoryReader() []Category {
	db := BddOpener()
	result, err1 := db.Query("SELECT CategoryID, CategoryName FROM Category")
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return nil
	}
	defer result.Close()
	var Categories []Category
	for result.Next() {
		var categoryName string
		var categoryID int
		result.Scan(&categoryID, &categoryName)
		var category = Category{categoryID, categoryName}
		Categories = append(Categories, category)
	}
	result.Close()
	db.Close()
	return Categories
}

// It returns the number of rows in the Post table
func HowManyRow() int {
	db := BddOpener()
	result, err1 := db.Query("SELECT COUNT(*) FROM Post")
	if err1 != nil {
		fmt.Println("ratio, ", err1)
		return 0
	}
	defer result.Close()
	var count int
	for result.Next() {
		result.Scan(&count)
	}
	result.Close()
	db.Close()
	return count
}
