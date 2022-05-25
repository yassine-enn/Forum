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
					// cookie := exe.CookieGenerator(r.FormValue("loginUsername"))
					// http.SetCookie(w, &cookie)
					session, _ := store.Get(r, "session")
					session.Values["userID"] = r.FormValue("loginUsername")
					session.Save(r, w)

				}
			} else {
				fmt.Println("pessi fraude finito")
				exe.PostTopic(r.FormValue("postContent"), r.FormValue("postTitle"))
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
		}
<<<<<<< HEAD
		// cookie, _ := r.Cookie(r.FormValue("loginUsername"))
		// user := http.Cookie{
		// 	Name: "Username", Value: r.FormValue("loginUsername"),
		// }
		// http.SetCookie(w, &user)
		// fmt.Println("user:", user)
		// fmt.Println(w, cookie)
		// data := Page{isLog}
=======
>>>>>>> 8c272f2fd2db5fe26795196fd3db1d65153d1f77
		fmt.Println(exe.PostDataReader())
		data := Page{isLog, exe.PostDataReader()}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})
<<<<<<< HEAD
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		_, ok := session.Values["userID"]
		if !ok {
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}
		data := Page{isLog, exe.PostDataReader()}
		tmpl.ExecuteTemplate(w, "index", data)
	})
=======
>>>>>>> 8c272f2fd2db5fe26795196fd3db1d65153d1f77
	fileServer := http.FileServer(http.Dir("./template/"))
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}
