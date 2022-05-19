package execute

import (
	"database/sql"
	"fmt"
)

func login(username string, password string) {
	db, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return
	}
	result, err1 := db.Query("SELECT Username, Password FROM User WHERE Username = ?", username)
	if err1 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err1)
		return
	}
	var Username string
	var Password string
	for result.Next() {
		result.Scan(&Username, &Password)
		fmt.Println(Username, Password)
	}
	result.Close()
	db.Close()

}
