package ImportFunction

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Login(username string, password string) bool {
	db, err := sql.Open("sqlite3", "./forumdb")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return false
	}
	result, err1 := db.Prepare("SELECT Username, PasswordHash FROM User WHERE Username = ?")
	if err1 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err1)
		return false
	}
	login, err2 := result.Query(username)
	if err2 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err2)
		return false
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
	if username == "" || email == "" || password == "" {
		fmt.Println("Veuillez remplir tous les champs")
		return
	}
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

	passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println("Erreur lors de la génération du hash", err)
		return
	}
	_, queryErr := statement.Exec(username, email, passwordHash)
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
