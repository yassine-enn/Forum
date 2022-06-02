package ImportFunction

import (
	"database/sql"
	"fmt"
)

// func LikePost(postID int, isLog bool, hasLiked bool, hasDiskliked bool) {
// 	if hasLiked {
// 		return
// 	}
// 	if hasDiskliked {
// 		hasDiskliked = false
// 	}
// 	if !isLog {
// 		return
// 	}
// 	db, err := sql.Open("sqlite3", "./forumdb")
// 	if err != nil {
// 		fmt.Println("Error when opening the DB:", err)
// 		return
// 	}
// 	defer db.Close()
// 	statement, prepareErr := db.Prepare("UPDATE Post SET likeCounter = likeCounter + 1 WHERE PostID = ?")
// 	if prepareErr != nil {
// 		fmt.Println("Error when preparing the statement:", prepareErr)
// 		return
// 	}
// 	_, queryErr := statement.Exec(postID)
// 	if queryErr != nil {
// 		fmt.Println("Error when querying the statement:", queryErr)
// 		return
// 	}
// 	db.Close()
// }

// func DislikePost(postID int, isLog bool, hasDisliked bool, hasLiked bool) {
// 	if hasDisliked {
// 		return
// 	}
// 	if hasLiked {
// 		hasLiked = false
// 	}
// 	if !isLog {
// 		return
// 	}
// 	db, err := sql.Open("sqlite3", "./forumdb")
// 	if err != nil {
// 		fmt.Println("Error when opening the DB:", err)
// 		return
// 	}
// 	// defer db.Close()
// 	statement, prepareErr := db.Prepare("UPDATE Post SET likeCounter = likeCounter - 1 WHERE PostID = ?")
// 	if prepareErr != nil {
// 		fmt.Println("Error when preparing the statement:", prepareErr)
// 		return
// 	}
// 	_, queryErr := statement.Exec(postID)
// 	if queryErr != nil {
// 		fmt.Println("Error when querying the statement:", queryErr)
// 		return
// 	}
// 	db.Close()
// }

// func hasAlreadyLiked(postID int, isLog bool) bool {
// 	if !isLog {
// 		return false
// 	}
// 	db, err := sql.Open("sqlite3", "./databaseForum")
// 	if err != nil {
// 		fmt.Println("Error when opening the DB:", err)
// 		return false
// 	}
// 	defer db.Close()
// 	statement, prepareErr := db.Prepare("SELECT likeCounter FROM Post WHERE PostID = ?")
// 	if prepareErr != nil {
// 		fmt.Println("Error when preparing the statement:", prepareErr)
// 		return false
// 	}
// 	result, queryErr := statement.Query(postID)
// 	if queryErr != nil {
// 		fmt.Println("Error when querying the statement:", queryErr)
// 		return false
// 	}
// 	defer result.Close()
// 	var likeCounter int
// 	for result.Next() {
// 		result.Scan(&likeCounter)
// 	}
// 	if likeCounter > 0 {
// 		return true
// 	}
// 	return false
// }

// func hasAlreadyLiked(username string, postID int) bool {

// }

// func hasAlreadyDisliked(username string, postID int) bool {

// }

func LikePostDb(postID int, username string, isLog bool) {
	if !isLog {
		return
	}
	if !HasAlreadyLiked(username, postID) {
		db, err := sql.Open("sqlite3", "./databaseForum")
		if err != nil {
			fmt.Println("Error when opening the DB:", err)
			return
		}
		defer db.Close()
		statement, prepareErr := db.Prepare("INSERT INTO Likes (Username, LikedPostID) VALUES (?, ?)")
		if prepareErr != nil {
			fmt.Println("Error when preparing the statement:", prepareErr)
			return
		}
		_, queryErr := statement.Exec(username, postID)
		if queryErr != nil {
			fmt.Println("Error when querying the statement:", queryErr)
			return
		}
		db.Close()
	} else {
		db, _ := sql.Open("sqlite3", "./databaseForum")
		statement, _ := db.Prepare("DELETE FROM Likes WHERE Username = ? AND LikedPostID = ?")
		_, _ = statement.Exec(username, postID)
		db.Close()
	}
}

func HasAlreadyLiked(username string, postID int) bool {
	db, err := sql.Open("sqlite3", "./databaseForum")
	if err != nil {
		fmt.Println("Error when opening the DB:", err)
		return false
	}
	defer db.Close()
	statement, prepareErr := db.Prepare("SELECT Username FROM Likes WHERE Username = ? AND LikedPostID = ?")
	if prepareErr != nil {
		fmt.Println("Error when preparing the statement:", prepareErr)
		return false
	}
	result, queryErr := statement.Query(username, postID)
	if queryErr != nil {
		fmt.Println("Error when querying the statement:", queryErr)
		return false
	}
	defer result.Close()
	var usernameResult string
	for result.Next() {
		result.Scan(&usernameResult)
	}
	if usernameResult == username {
		return true
	}
	return false
}

// func DislikePostDB(postID int, username string, isLog bool) {
// 	if !isLog {
// 		return
// 	}
// 	if HasAlreadyLiked(username, postID) {
// 		db, _ := sql.Open("sqlite3", "./databaseForum")
// 		statement, _ := db.Prepare("DELETE FROM Likes WHERE Username = ? AND LikedPostID = ?")
// 		if !HasAlreadyDisliked(username, postID) {
// 			db, err := sql.Open("sqlite3", "./databaseForum")
// 			if err != nil {
// 				fmt.Println("Error when opening the DB:", err)
// 				return
// 			}
// 			defer db.Close()
// 			statement, prepareErr := db.Prepare("INSERT INTO Dislikes (Username, DislikedPost) VALUES (?, ?)")
// 			if prepareErr != nil {
// 				fmt.Println("Error when preparing the statement:", prepareErr)
// 				return
// 			}
// 			_, queryErr := statement.Exec(username, postID)
// 			if queryErr != nil {
// 				fmt.Println("Error when querying the statement:", queryErr)
// 				return
// 			}
// 			db.Close()
// 		} else {
// 			db, _ := sql.Open("sqlite3", "./databaseForum")
// 			statement, _ := db.Prepare("DELETE FROM Dislikes WHERE Username = ? AND DislikedPostID = ?")
// 			_, _ = statement.Exec(username, postID)
// 			db.Close()
// 		}
// 	}
// }

func HasAlreadyDisliked(username string, postID int) bool {
	db, err := sql.Open("sqlite3", "./databaseForum")
	if err != nil {
		fmt.Println("Error when opening the DB:", err)
		return false
	}
	defer db.Close()
	statement, prepareErr := db.Prepare("SELECT Username FROM Dislikes WHERE Username = ? AND DislikedPostID = ?")
	if prepareErr != nil {
		fmt.Println("Error when preparing the statement:", prepareErr)
		return false
	}
	result, queryErr := statement.Query(username, postID)
	if queryErr != nil {
		fmt.Println("Error when querying the statement:", queryErr)
		return false
	}
	defer result.Close()
	var usernameResult string
	for result.Next() {
		result.Scan(&usernameResult)
	}
	if usernameResult == username {
		return true
	}
	return false
}
