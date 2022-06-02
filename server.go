package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	exe "ImportFunction/ImportFunction"

	"github.com/google/uuid"
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
var wichPost int
var username string

func main() {
	var sessionToken string
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	http.HandleFunc("/", Er404)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if r.FormValue("signup_button") == "LOG IN" {
				isLog, isCorrectPwd = exe.Login(r.FormValue("loginUsername"), r.FormValue("loginPassword"), sessions[sessionToken])
				fmt.Println("isLog", isLog, "isCorrectPwd", isCorrectPwd)
				if isCorrectPwd {
					sessionToken = uuid.NewString()
					expiresAt := time.Now().Add(3600 * time.Second)
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
						username = r.FormValue("loginUsername")
					}
				}
			} else if r.FormValue("signup_button") == "SIGN UP" {
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
			fmt.Println("wichPost", wichPost)
		}
		fmt.Println("islogF", isLog)
		data := Page{isLog, exe.PostDataReader("PostID > 0")}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})
	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/dislike", dislikeHandler)

	http.HandleFunc("/postCreator", func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session_token")
		sessionToken := c.Value
		userSession, _ := sessions[sessionToken]
		if userSession.IsExpired() {
			delete(sessions, sessionToken)
			return
		}
		fmt.Println(userSession.Username)
		if r.Method == "POST" {
			exe.PostTopic(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_category"), userSession.Username)
		}
		data := Page{isLog, exe.PostDataReader("PostID > 0")}
		tmpl.ExecuteTemplate(w, "postcreator", data)
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session_token")
		sessionToken := c.Value
		userSession, _ := sessions[sessionToken]
		if userSession.IsExpired() {
			delete(sessions, sessionToken)
			return
		}
		isLog = true
		wichPost, _ = strconv.Atoi(r.FormValue("post_id"))
		data := Page{isLog, exe.PostDataReader("PostID = " + strconv.Itoa(wichPost))}
		tmpl.ExecuteTemplate(w, "post_page", data)
	})
	fileServer := http.FileServer(http.Dir("./template/"))
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession, _ := sessions[sessionToken]
	if userSession.IsExpired() {
		delete(sessions, sessionToken)
		return
	}
	isLog = true
	fmt.Println("like", r.FormValue("post_id_like"))
	postIdLike, _ := strconv.Atoi(r.FormValue("post_id_like"))
	exe.GetLikes(postIdLike)
	exe.LikePostDb(postIdLike, username, isLog)
	exe.GetLikes(postIdLike)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func dislikeHandler(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	sessionToken := c.Value
	userSession, _ := sessions[sessionToken]
	if userSession.IsExpired() {
		delete(sessions, sessionToken)
		return
	}
	isLog = true
	fmt.Println("dislike", r.FormValue("post_id_dislike"), username)
	postIdDislike, _ := strconv.Atoi(r.FormValue("post_id_dislike"))
	exe.DislikePostDB(postIdDislike, username, isLog)
	exe.GetDislikes(postIdDislike)
	exe.GetLikes(postIdDislike)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
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

func Er404(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("template/vues/err404.html")).ExecuteTemplate(w, "err404.html", nil) //opening of the er404.html page
}
