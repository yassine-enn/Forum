package main

import (
	"fmt"
	"html/template"
	"net/http"

	exe "ImportFunction/ImportFunction"
)

func main() {
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if r.FormValue("signup_button") == "LOG IN" {
				// exe.Login(r.FormValue("loginUsername"), r.FormValue("loginPassword"))
			} else {
				fmt.Println("pessi fraude finito")
				exe.PostTopic(r.FormValue("postContent"), r.FormValue("postTitle"))
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
		}
		tmpl.ExecuteTemplate(w, "acceuil", nil)
	})
	fileServer := http.FileServer(http.Dir("./template/"))
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
