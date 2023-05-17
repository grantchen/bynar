package http_handler

import (
	"log"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

// TODO: get throug request
var (
	ModuleID  int = 4
	AccountID int = 123456
)

func NewHTTPHandler(appConfig config.AppConfig) *handler.HTTPTreeGridHandler {
	dbConnString := appConfig.GetDBConnection()

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
			return transferService.HandleUpload(req, AccountID)
		},
	}
	return handler
}
