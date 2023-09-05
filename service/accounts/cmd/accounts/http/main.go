package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/http_handler"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
)

// use for test only this module without permission
func main() {
	connString := "root:123456@tcp(localhost:3306)/accounts_management"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	handler := http_handler.NewHTTPHandler(db)

	http.HandleFunc("/signup", handler.Signup)
	http.HandleFunc("/signin", handler.Signin)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
