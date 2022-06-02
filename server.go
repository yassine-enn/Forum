package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"time"

	exe "ImportFunction/ImportFunction"

	"github.com/google/uuid"
)

var sessions = map[string]exe.Session{}

type Page struct {
	IsLoged     bool
	Post        []exe.Post
	Comment     []exe.Comment
	Categories  []exe.Category
	HowManyPage int
	WhichPage   int
}

var isLog bool
var isCorrectPwd bool
var wichPost int
var username string
var postToShow []exe.Post
var paginValue = 1
var whatPage = 1

func main() {
	var sessionToken string
	tmpl, err := template.ParseGlob("./template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}
	// It's a function that redirects to the er404.html page if the user enters a wrong url.
	http.HandleFunc("/", Er404)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// This is the code that allows the user to log in.
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
			} // This is the code that allows the user to sign up.
		}
		// It's the code that allows the user to go to the next page or the previous page.
		if r.FormValue("next_page") == "next" {
			paginValue += 10
			whatPage += 1
		} else if r.FormValue("previous_page") == "previous" {
			paginValue -= 10
			whatPage -= 1
		}
		whichPage := strconv.Itoa(paginValue)
		maxPage := math.Ceil(float64(exe.HowManyRow()) / 10)
		postToShow = exe.PostDataReader("PostID > 0", `LIMIT '10' OFFSET `+whichPage)
		// It's the code that allows the user to search for a post.
		if r.FormValue("search") != "" {
			postToShow = exe.PostDataReader("PostText LIKE '%"+r.FormValue("search")+"%' OR PostTitle LIKE '%"+r.FormValue("search")+"%' OR PostCategory LIKE '%"+r.FormValue("search")+"%' OR PostAuthor LIKE '%"+r.FormValue("search")+"%'", "")

		} else if r.FormValue("category_filter") != "" {
			postToShow = exe.PostDataReader("PostCategory = '"+r.FormValue("category_filter")+"'", "")
		} // It's the code that allows the user to filter the posts by category.
		data := Page{isLog, postToShow, nil, exe.CategoryReader(), int(maxPage), whatPage}
		tmpl.ExecuteTemplate(w, "acceuil", data)
	})
	http.HandleFunc("/like", likeHandler)
	http.HandleFunc("/dislike", dislikeHandler)

	http.HandleFunc("/postCreator", func(w http.ResponseWriter, r *http.Request) {
		// It's a code that allows the user to stay logged in, by checking if a session is opened.
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
		data := Page{isLog, nil, nil, nil, 0, 0}
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
		data := Page{isLog, exe.PostDataReader("PostID = "+strconv.Itoa(wichPost), "10"), exe.CommentDataReader("CommentSource = " + strconv.Itoa(wichPost)), exe.CategoryReader(), 0, 0}
		tmpl.ExecuteTemplate(w, "post_page", data)
	})

	fileServer := http.FileServer(http.Dir("./template/"))
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/template/", http.StripPrefix("/template/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
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
