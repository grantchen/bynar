package main

import (
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/config"
	http_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

func main() {

	appConfig := config.NewLocalConfig()
	dbConnString := appConfig.GetDBConnection()

	db, err := sql_db.NewConnection(dbConnString)
	if err != nil {
		log.Panic(err)
	}

	transferRepository := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(transferRepository)
	transferHandler := http_handler.NewTransferHttpHandler(transferService)
	http.HandleFunc("/data", transferHandler.HandleGetPageCount)
	http.HandleFunc("/page", transferHandler.HandleGetPageTransferData)
	log.Println("start server at 8081!")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
