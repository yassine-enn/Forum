package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		var UserID int
		var userName string
		var userEmail string
		var userPassword string
		var userAvatar string
		if r.Method == "POST" {
			UserID = 5
			userName = r.FormValue("signupUsername")
			userEmail = r.FormValue("signupEmail")
			userPassword = r.FormValue("signupPassword")
			userAvatar = "sasuke"
			db, err := sql.Open("sqlite3", "./Forum.db")
			if err != nil {
				fmt.Println("Echec de l'ouverture de la base", err)
				return
			}
			statement, prepareErr := db.Prepare("INSERT INTO User (UserID, Username, Email, PasswordHash, Avatar) VALUES (?,?,?,?,?)")
			if prepareErr != nil {
				fmt.Println("La préparation de la requête a échoué", prepareErr)
				return
			}
			_, queryErr := statement.Exec(UserID, userName, userEmail, userPassword, userAvatar)
			if queryErr != nil {
				fmt.Println("Une erreur est survenue durant la requête", queryErr)
				return
			}
			statement.Close()
			db.Close()
			fmt.Println("username:", r.FormValue("signupUsername"))
			fmt.Println("email:", r.FormValue("signupEmail"))
			fmt.Println("password:", r.FormValue("signupPassword"))
			fmt.Println("password_confirmation:", r.FormValue("signupPasswordConfirm"))
		}
		tmpl.ExecuteTemplate(w, "acceuil", nil)
	})
	fileServer := http.FileServer(http.Dir("./template/"))
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
