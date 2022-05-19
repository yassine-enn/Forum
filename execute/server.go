package forum_server

import (
	"fmt"
	"html/template"
	"net/http"
)

func server() {
	tmpl, err := template.ParseGlob("/template/vues/*.html")
	if err != nil {
		fmt.Println("Template loading Error:", err)
		return
	}

	http.HandleFunc("../template/vues/accueil.html", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "accueil.html", nil)
	})
	fileServer := http.FileServer(http.Dir("../template/style-forum"))
	http.Handle("/style-forum/", http.StripPrefix("/style-forum/", fileServer))
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
