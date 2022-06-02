package ImportFunction

import (
	"database/sql"
	"fmt"
)

func LikePostDb(postID int, username string, isLog bool) {
	if !isLog {
		return
	}
	if HasAlreadyDisliked(username, postID) {
		db, _ := sql.Open("sqlite3", "./databaseForum")
		statement, _ := db.Prepare("DELETE FROM Dislikes WHERE Username = ? AND DislikedPostID = ?")
		statement.Exec(username, postID)
		db.Close()
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
	db := BddOpener()
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

func DislikePostDB(postID int, username string, isLog bool) {
	if !isLog {
		return
	}
	if HasAlreadyLiked(username, postID) {
		db := BddOpener()
		statement, prepareErr := db.Prepare("DELETE FROM Likes WHERE Username = ? AND LikedPostID = ?")
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
	}
	if !HasAlreadyDisliked(username, postID) {
		db := BddOpener()
		defer db.Close()
		statement, prepareErr := db.Prepare("INSERT INTO Dislikes (Username, DislikedPostID) VALUES (?, ?)")
		if prepareErr != nil {
			fmt.Println("Error when preparing the statement:", prepareErr)
			return
		}
		_, queryErr := statement.Exec(username, postID)
		if queryErr != nil {
			fmt.Println("Error when querying the statement:", queryErr)
			return
		} else {
			fmt.Println("Dislike added")
		}
		db.Close()
	} else {
		db, _ := sql.Open("sqlite3", "./databaseForum")
		statement, _ := db.Prepare("DELETE FROM Dislikes WHERE Username = ? AND DislikedPostID = ?")
		_, _ = statement.Exec(username, postID)
		db.Close()
	}
}

func HasAlreadyDisliked(username string, postID int) bool {
	db := BddOpener()
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
		fmt.Println("Dislike found")
		return true
	}
	return false
}

func GetLikes(postID int) {
	db := BddOpener()
	defer db.Close()
	statement, prepareErr := db.Prepare("SELECT COUNT(*) FROM Likes WHERE LikedPostID = ?")
	if prepareErr != nil {
		fmt.Println("Error when preparing the statement:", prepareErr)
		return
	}
	result, queryErr := statement.Query(postID)
	if queryErr != nil {
		fmt.Println("Error when querying the statement:", queryErr)
		return
	}
	defer result.Close()
	var likes int
	for result.Next() {
		result.Scan(&likes)
	}
	fmt.Println(likes, "likes")
	db, _ = sql.Open("sqlite3", "./databaseForum")
	statement, _ = db.Prepare("UPDATE Post SET likeCounter = ? WHERE PostID = ?")
	_, _ = statement.Exec(likes-GetDislikes(postID), postID)
	db.Close()
}

func GetDislikes(postID int) int {
	db := BddOpener()
	defer db.Close()
	statement, prepareErr := db.Prepare("SELECT COUNT(*) FROM Dislikes WHERE DislikedPostID = ?")
	if prepareErr != nil {
		fmt.Println("Error when preparing the statement:", prepareErr)
		return 0
	}
	result, queryErr := statement.Query(postID)
	if queryErr != nil {
		fmt.Println("Error when querying the statement:", queryErr)
		return 0
	}
	defer result.Close()
	var dislikes int
	for result.Next() {
		result.Scan(&dislikes)
	}
	fmt.Println(dislikes, "dislikes")
	return dislikes
}
