package main

import (
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

var accountID = 11111

func main() {

	// appConfig := config.NewLocalConfig()
	// dbConnString := appConfig.GetDBConnection()
	dbConnString := "root:123456@tcp(localhost:3306)/bynar"

	db, err := sql_db.NewConnection(dbConnString)
	if err != nil {
		log.Panic(err)
	}

	documentRepository := repository.NewDocumentRepository(db)
	inventoryRepository := repository.NewInventoryRepository(db)
	transferRepository := repository.NewTransferRepository(db)
	userRepository := repository.NewUserRepository(db)
	workflowRepository := repository.NewWorkflowRepository()

	transferService := service.NewTransferService(
		db,
		userRepository,
		workflowRepository,
		transferRepository,
		inventoryRepository,
		documentRepository,
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackGetPageCountFunc: transferService.GetPagesCount,
		CallbackGetPageDataFunc:  transferService.GetTransfersPageData,
		CallbackUploadDataFunc: func(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
			return transferService.HandleUpload(req, accountID)
		},
	}

	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
