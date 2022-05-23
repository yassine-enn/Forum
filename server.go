package main

import (
	"fmt"
	"html/template"
	"net/http"

	exe "ImportFunction/ImportFunction"
)

type Page struct {
	IsLoged bool
}

func main() {
	var isLog bool
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if r.FormValue("signup_button") == "LOG IN" {
				isLog = exe.Login(r.FormValue("loginUsername"), r.FormValue("loginPassword"))
				if isLog {
					cookie := exe.CookieGenerator(r.FormValue("loginUsername"))
					http.SetCookie(w, &cookie)
				}
			} else {
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
			if r.FormValue("logout") == "logout" {
				isLog = false
			} else {
				fmt.Println("No logout")
			}
		}
		data := Page{isLog}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})
	fileServer := http.FileServer(http.Dir("./template/"))
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
