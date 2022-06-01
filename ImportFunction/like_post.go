package ImportFunction

import (
	"database/sql"
	"fmt"
)

func LikePost(postID int) {
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Error when opening the BDD:", err)
		return
	}
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
}
