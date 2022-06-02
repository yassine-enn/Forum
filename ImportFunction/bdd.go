package ImportFunction

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func BddOpener() *sql.DB {
	db, err := sql.Open("sqlite3", "./dbforum")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return nil
	}
	return db
}
