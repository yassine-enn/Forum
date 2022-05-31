package ImportFunction

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Login(username string, password string, session Session) (bool, bool) {
	db, err := sql.Open("sqlite3", "./ALED")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base")
		return false, false
	}
	result, err1 := db.Prepare("SELECT Username, PasswordHash FROM User WHERE Username = ?")
	if err1 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err1)
		return false, false
	}
	login, err2 := result.Query(username)
	if err2 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err2)
		return false, false
	}
	var UsernameFromDataBase string
	var PasswordFromDataBase string
	for login.Next() {
		login.Scan(&UsernameFromDataBase, &PasswordFromDataBase)
		if err := bcrypt.CompareHashAndPassword([]byte(PasswordFromDataBase), []byte(password)); err != nil {
			fmt.Println("wrong password")
			return false, false
		} else {
			fmt.Println("password was correct")
			if !session.IsExpired() {
				return true, true
			} else {
				return false, true
			}
		}
	}
	result.Close()
	db.Close()
	return false, false
}

func Signup(username string, email string, password string) string {
	if username == "" || email == "" || password == "" {
		fmt.Println("Veuillez remplir tous les champs")
		return "Veuillez remplir tous les champs"
	}
	if AlreadyExist(username, email) {
		return "Ce nom d'utilisateur ou cet email existe déjà"
	}
	db, err := sql.Open("sqlite3", "./ALED")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base", err)
		return ""
	}
	statement, prepareErr := db.Prepare("INSERT INTO User (Username, Email, PasswordHash) VALUES (?,?,?)")
	if prepareErr != nil {
		fmt.Println("La préparation de la requête a échoué", prepareErr)
		return ""
	}

	passwordHash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		fmt.Println("Erreur lors de la génération du hash", err)
		return ""
	}
	_, queryErr := statement.Exec(username, email, passwordHash)
	if queryErr != nil {
		fmt.Println("Une erreur est survenue durant la requête", queryErr)
		return ""
	}
	statement.Close()
	db.Close()
	fmt.Println("username:", username)
	fmt.Println("email:", email)
	fmt.Println("password:", password)
	return "Votre compte a bien été créé"
}

func AlreadyExist(username string, email string) bool {
	db, err := sql.Open("sqlite3", "./ALED")
	if err != nil {
		fmt.Println("Echec de l'ouverture de la base", err)
		return false
	}
	result, err1 := db.Prepare("SELECT Username, Email FROM User WHERE Username = ? OR Email = ?")
	if err1 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err1)
		return false
	}
	login, err2 := result.Query(username, email)
	if err2 != nil {
		fmt.Println("erreur lors de la recherche dans la base de donnée", err2)
		return false
	}
	var UsernameFromDataBase string
	var EmailFromDataBase string
	for login.Next() {
		login.Scan(&UsernameFromDataBase, &EmailFromDataBase)
		if UsernameFromDataBase == username || EmailFromDataBase == email {
			fmt.Println("username already exist or email already exist")
			return true
		}
	}
	result.Close()
	db.Close()
	return false
}
