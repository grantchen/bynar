package main

import (
	"log"
	"net/http"
	"os"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/task2/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/task2/internal/service"
)

var accountID = 11111

func main() {
	dbConnString := os.Getenv("DB_CONN")
	db, err := sql_db.NewConnection(dbConnString)
	if err != nil {
		log.Panic(err)
	}
	documentRepository := repository.NewDocumentRepository(db)
	inventoryRepository := repository.NewInventoryRepository(db)
	transferRepository := repository.NewTransferRepository(db)
	userRepository := repository.NewUserRepository(db)
	workflowRepository := repository.NewWorkflowRepository()

	uploadService := service.NewUploadService(
		db,
		userRepository,
		workflowRepository,
		transferRepository,
		inventoryRepository,
		documentRepository,
	)
	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: func(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
			return uploadService.Handle(req, accountID)
		},
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)

}
