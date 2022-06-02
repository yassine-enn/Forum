package ImportFunction

import (
	"fmt"
)

func LikePost(postID int) {
	db := BddOpener()
	defer db.Close()
	statement, prepareErr := db.Prepare("UPDATE Post SET likeCounter = likeCounter + 1 WHERE PostID = ?")
	if prepareErr != nil {
		fmt.Println("Error when preparing the statement:", prepareErr)
		return
	}
	_, queryErr := statement.Exec(postID)
	if queryErr != nil {
		fmt.Println("Error when querying the statement:", queryErr)
		return
	}
	db.Close()
}

func DislikePost(postID int) {
	db := BddOpener()
	// defer db.Close()
	statement, prepareErr := db.Prepare("UPDATE Post SET likeCounter = likeCounter - 1 WHERE PostID = ?")
	if prepareErr != nil {
		fmt.Println("Error when preparing the statement:", prepareErr)
		return
	}
	_, queryErr := statement.Exec(postID)
	if queryErr != nil {
		fmt.Println("Error when querying the statement:", queryErr)
		return
	}
	db.Close()
}
