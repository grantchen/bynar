package main

import (
	"log"
	"net/http"

	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

var accountID = 6

// only for testing
func main() {
	// appConfig := config.NewLocalConfig()
	// dbConnString := appConfig.GetDBConnection()
	dbConnString := "root:123456@tcp(localhost:3306)/46542255-9d45-49d5-939d-84bc55b1a938"

	db, err := sqldb.NewConnection(dbConnString)
	if err != nil {
		log.Panic(err)
	}

	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"transfers",
		"transfer_lines",
		repository.TransferFieldNames,
		repository.TransferLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "document_id",
			QueryParent:              repository.QueryParent,
			QueryParentCount:         repository.QueryParentCount,
			QueryParentJoins:         repository.QueryParentJoins,
			QueryChild:               repository.QueryChild,
			QueryChildCount:          repository.QueryChildCount,
			QueryChildJoins:          repository.QueryChildJoins,
			QueryChildSuggestion:     repository.QueryChildSuggestion,
			ChildJoinFieldWithParent: "parent_id",
			ParentIdField:            "id",
		},
	)

	grTransferRepositoryWithChild := treegrid.NewGridRepository(db,
		"transfers",
		"transfer_lines",
		repository.TransferFieldNames,
		repository.TransferLineFieldNames,
	)
	documentRepository := repository.NewDocumentRepository(db)
	inventoryRepository := repository.NewInventoryRepository(db)
	transferRepository := repository.NewTransferRepository(db, "en")
	userRepository := repository.NewUserRepository(db)

	transferService := service.NewTransferService(
		db,
		transferRepository,
		gridRowDataRepositoryWithChild,
	)

	uploadService := service.NewUploadService(
		db,
		grTransferRepositoryWithChild,
		userRepository,
		transferRepository,
		inventoryRepository,
		documentRepository,
		accountID,
		"en",
	)
	h := &handler.HTTPTreeGridHandler{
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := transferService.GetPageCount(tr)
			return float64(count), err
		},
		CallbackGetPageDataFunc: transferService.GetPageData,
		CallbackUploadDataFunc: func(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
			return uploadService.Handle(req)
		},
	}

	http.HandleFunc("/data", h.HTTPHandleGetPageCount)
	http.HandleFunc("/page", h.HTTPHandleGetPageData)
	http.HandleFunc("/upload", h.HTTPHandleUpload)
	log.Println("start server at 8081!")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
