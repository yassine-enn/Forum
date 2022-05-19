package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		r.FormValue("username")
		r.FormValue("email")
		r.FormValue("password")
		r.FormValue("password_confirmation")
		fmt.Println("username:", r.FormValue("username"))
		fmt.Println("email:", r.FormValue("email"))
		fmt.Println("password:", r.FormValue("password"))
		fmt.Println("password_confirmation:", r.FormValue("password_confirmation"))

		tmpl.ExecuteTemplate(w, "acceuil", nil)
	})
	fileServer := http.FileServer(http.Dir("./template/"))
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
