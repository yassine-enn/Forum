package main

import (
	"fmt"
	"html/template"
	"net/http"

	exe "ImportFunction/ImportFunction"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

type Page struct {
	IsLoged bool
	Post    []exe.Post
}

var store = sessions.NewCookieStore([]byte("super-secret"))

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
					session, _ := store.Get(r, "session")
					session.Values["userID"] = r.FormValue("loginUsername")
					session.Save(r, w)
				}
			} else {
				fmt.Println("pessi fraude finito")
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
		}
		fmt.Println(exe.PostDataReader())
		data := Page{isLog, exe.PostDataReader()}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})

	http.HandleFunc("/postCreator", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "postcreator", nil)
		if r.Method == "POST" {
			// exe.PostTopic(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_category"), image)
			fmt.Println(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_category"))
		}
	})

	fileServer := http.FileServer(http.Dir("./template/"))
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
