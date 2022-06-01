package ImportFunction

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

<<<<<<< HEAD
func BddReader(wichTable string, condition string) []Post {
	var posts []Post
	db, err := sql.Open("sqlite3", "./forumdb")
=======
func BddOpener() *sql.DB {
	db, err := sql.Open("sqlite3", "./ALED")
>>>>>>> 35b068193f0429612d9fe07d79731cbd1afbd2e3
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	return db
}
