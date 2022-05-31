package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	exe "ImportFunction/ImportFunction"

	"github.com/google/uuid"
	"github.com/gorilla/context"
)

var sessions = map[string]exe.Session{}

// func (s exe.Session) isExpired() bool {
// 	return s.Expiry.Before(time.Now())
// }

type Page struct {
	IsLoged bool
	Post    []exe.Post
}

var isLog bool
var isCorrectPwd bool

var data Page

func main() {
	var sessionToken string
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if r.FormValue("signup_button") == "LOG IN" {
				isLog, isCorrectPwd = exe.Login(r.FormValue("loginUsername"), r.FormValue("loginPassword"), sessions[sessionToken])
				fmt.Println("isLog", isLog, "isCorrectPwd", isCorrectPwd)
				if isCorrectPwd {
					sessionToken = uuid.NewString()
					expiresAt := time.Now().Add(120 * time.Second)
					// Set the token in the session map, along with the session information
					sessions[sessionToken] = exe.Session{
						Username: r.FormValue("loginUsername"),
						Expiry:   expiresAt,
					}
					isLog, _ = exe.Login(r.FormValue("loginUsername"), r.FormValue("loginPassword"), sessions[sessionToken])
					fmt.Println(isLog)
					// Finally, we set the client cookie for "session_token" as the session token we just generated
					// we also set an expiry time of 120 seconds
					if isLog {
						http.SetCookie(w, &http.Cookie{
							Name:    "session_token",
							Value:   sessionToken,
							Expires: expiresAt,
						})
						fmt.Println("2", sessionToken)
						fmt.Println("s", sessions)
					}
				}
				data = Page{isLog, exe.PostDataReader()}
				tmpl.ExecuteTemplate(w, "acceuil", data)
			} else {
				fmt.Println("pessi fraude finito")
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
		}
		fmt.Println("islogF", isLog)
		data := Page{isLog, exe.PostDataReader()}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		data := Page{isLog, exe.PostDataReader()}
		fmt.Println("bb", data.IsLoged)
		tmpl.ExecuteTemplate(w, "index", data)
	})
	http.HandleFunc("/postCreator", func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session_token")
		sessionToken := c.Value
		userSession, _ := sessions[sessionToken]
		if userSession.IsExpired() {
			delete(sessions, sessionToken)
			return
		}
		isLog = true
		tmpl.ExecuteTemplate(w, "postcreator", data)
		if r.Method == "POST" {
			// exe.PostTopic(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_category"), image)
			fmt.Println(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_category"))
		}
	})

	fileServer := http.FileServer(http.Dir("./template/"))
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// tpl, _ := template.ParseFiles("./template/vues/logout.html")
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	delete(sessions, sessionToken)
	fmt.Println("s1", sessions)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	isLog = false
	redirect := "/home"
	http.Redirect(w, r, redirect, http.StatusFound)
}
