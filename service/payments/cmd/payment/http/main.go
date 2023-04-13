package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/config"
	http_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/handler/http"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
)

// TODO: get throug request
var (
	ModuleID  int = 8
	AccountID int = 123456
)

func main() {

	appConfig := config.NewLocalConfig()
	dbConnString := appConfig.GetDBConnection()

	if _, err := sql_db.NewConnection(dbConnString); err != nil {
		log.Fatalln("new connection", err)
	}

	uploadHandler := &http_handler.UploadHandler{ModuleID: ModuleID, AccountID: AccountID}
	http.HandleFunc("/upload", uploadHandler.HandleUpload)
	http.HandleFunc("/test", uploadHandler.HandleUploadTest)
	log.Println("start server at 8081!")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
