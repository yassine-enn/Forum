package ImportFunction

import (
	"database/sql"
	"fmt"
)

func Login(username string, password string) {
	db, err := sql.Open("sqlite3", "./Forum.db")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return
	}
	result, err1 := db.Prepare("SELECT Username, Password FROM User WHERE Username = ?")
	if err1 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err1)
		return
	}
	login, err2 := result.Query(username)
	if err2 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err2)
		return
	}
	var Username string
	var Password string
	for login.Next() {
		login.Scan(&Username, &Password)
		fmt.Println(Username, Password)
	}
	result.Close()
	db.Close()
}

func Signup(username string, email string, password string) {
	db, err := sql.Open("sqlite3", "./ALED")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base", err)
		return
	}
	statement, prepareErr := db.Prepare("INSERT INTO User (Username, Email, PasswordHash) VALUES (?,?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return
	}
	_, queryErr := statement.Exec(username, email, password)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return
	}
	statement.Close()
	db.Close()
	fmt.Println("username:", username)
	fmt.Println("email:", email)
	fmt.Println("password:", password)
}
