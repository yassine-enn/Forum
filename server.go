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
	IsLoged    bool
	Post       []exe.Post
	Comment    []exe.Comment
	Categories []exe.Category
}

var isLog bool
var isCorrectPwd bool
var wichPost int
var username string
var postToShow []exe.Post

// var hasLiked = false
// var hasDisliked = false
// var likeNum = 0
// var dislikeNum = 0
// var data Page

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
						username = r.FormValue("loginUsername")
					}
				}
			} else if r.FormValue("signup_button") == "SIGN UP" {
				exe.Signup(r.FormValue("signupUsername"), r.FormValue("signupEmail"), r.FormValue("signupPassword"))
			}
			fmt.Println("wichPost", wichPost)
		}

		// fmt.Println("like", r.FormValue("post_id_like"))
		// postIdLike, _ := strconv.Atoi(r.FormValue("post_id_like"))
		// if r.FormValue("post_id_like") != "" {
		// 	exe.LikePost(postIdLike)
		// }
		// fmt.Println("dislike", r.FormValue("post_id_dislike"))
		// postIdDislike, _ := strconv.Atoi(r.FormValue("post_id_dislike"))
		// if r.FormValue("post_id_dislike") != "" {
		// 	exe.DislikePost(postIdDislike)
		// }
		postToShow = exe.PostDataReader("PostID > 0")
		if r.FormValue("search") != "" {
			postToShow = exe.PostDataReader("PostText LIKE '%" + r.FormValue("search") + "%' OR PostTitle LIKE '%" + r.FormValue("search") + "%'")
		}
		fmt.Println("islogF", isLog)
		data := Page{isLog, postToShow, nil, exe.CategoryReader()}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})
	http.HandleFunc("/like", likeHandler)
	// http.HandleFunc("/dislike", dislikeHandler)

	http.HandleFunc("/postCreator", func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie("session_token")
		if c == nil {
			redirect := "/home"
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		sessionToken := c.Value
		if sessionToken == "" {
			redirect := "/home"
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		userSession, _ := sessions[sessionToken]
		if userSession.IsExpired() {
			isLog = false
			redirect := "/home"
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		if r.Method == "POST" {
			if r.FormValue("post_input_new_category") != "" {
				exe.AddCategory(r.FormValue("post_input_new_category"))
				exe.PostTopic(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_new_category"), userSession.Username)
			} else {
				exe.PostTopic(r.FormValue("post_input_text"), r.FormValue("post_input_title"), r.FormValue("post_input_category"), userSession.Username)
			}
		}
		data := Page{isLog, exe.PostDataReader("PostID > 0"), nil, exe.CategoryReader()}
		tmpl.ExecuteTemplate(w, "postcreator", data)
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		wichPost, _ = strconv.Atoi(r.FormValue("post_id"))
		c, _ := r.Cookie("session_token")
		if c == nil {
			redirect := "/home"
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		sessionToken := c.Value
		if sessionToken == "" {
			redirect := "/home"
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		userSession := sessions[sessionToken]
		if userSession.IsExpired() {
			isLog = false
			redirect := "/home"
			http.Redirect(w, r, redirect, http.StatusFound)
			return
		}
		if r.Method == "POST" {
			wichPost, _ = strconv.Atoi(r.FormValue("post_id"))
			if r.FormValue("comment_content") != "" {
				exe.CommentTopic(r.FormValue("comment_content"), wichPost, userSession.Username)
			}
		}
		data := Page{isLog, exe.PostDataReader("PostID = " + strconv.Itoa(wichPost)), exe.CommentDataReader("CommentSource = " + strconv.Itoa(wichPost)), exe.CategoryReader()}
		tmpl.ExecuteTemplate(w, "post_page", data)
	})

	fileServer := http.FileServer(http.Dir("./template/"))
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("like", r.FormValue("post_id_like"))
	postIdLike, _ := strconv.Atoi(r.FormValue("post_id_like"))
	exe.LikePostDb(postIdLike, username, isLog)
	http.Redirect(w, r, "/home", http.StatusFound)
}

// func dislikeHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("dislike", r.FormValue("post_id_dislike"))
// 	postIdDislike, _ := strconv.Atoi(r.FormValue("post_id_dislike"))
// 	exe.DislikePostDB(postIdDislike, username, isLog)
// 	http.Redirect(w, r, "/home", http.StatusFound)
// }

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
