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

// each session contains the username of the user and the time at which it expires
// type session struct {
// 	username string
// 	expiry   time.Time
// }

// we'll use this method later to determine if the session has expired

type Page struct {
	IsLoged bool
	Post    []exe.Post
}

// var data Page

// var store = sessions.NewCookieStore([]byte("super-secret"))
var isLog bool
var isCorrectPwd bool

func main() {
	// var isLog bool
	// var isCorrectPwd bool
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
				data := Page{isLog, exe.PostDataReader()}
				tmpl.ExecuteTemplate(w, "acceuil", data)
			} else {
				fmt.Println("pessi fraude finito")
				exe.PostTopic(r.FormValue("postContent"), r.FormValue("postTitle"))
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
	fileServer := http.FileServer(http.Dir("./template/"))
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("./template/vues/logout.html")
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
	tpl.ExecuteTemplate(w, "logout.html", nil)

}
