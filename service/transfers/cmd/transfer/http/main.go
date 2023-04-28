package main

import (
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

func main() {

	// appConfig := config.NewLocalConfig()
	// dbConnString := appConfig.GetDBConnection()
	dbConnString := "root:123456@tcp(localhost:3306)/bynar"

	db, err := sql_db.NewConnection(dbConnString)
	if err != nil {
		log.Panic(err)
	}

	transferRepository := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(transferRepository)

	handler := &handler.HTTPTreeGridHandler{CallbackGetPageCountFunc: transferService.GetPagesCount,
		CallbackGetPageDataFunc: transferService.GetTransfersPageData}

	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
