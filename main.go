// package main

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/mattn/go-sqlite3"
// )

// func main() {
// 	db, err := sql.Open("sqlite3", "./Forum.db")
// 	if err != nil {
// 		fmt.Println("Echec de l'ouverture de la base")
// 		return
// 	}
// 	result, err1 := db.Query("SELECT Username, Email FROM User WHERE UserID > 0")
// 	if err1 != nil {
// 		fmt.Println("ratio, ", err1)
// 		return
// 	}
// 	var Username string
// 	var Email string
// 	for result.Next() {
// 		result.Scan(&Username, &Email)
// 		fmt.Println(Username, Email)
// 	}
// 	result.Close()
// 	db.Close()
// }
