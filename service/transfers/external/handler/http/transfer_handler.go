package http_handler

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
)

// TODO: get throug request
var (
	AccountID int = 123456
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {

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
	workflowRepository := repository.NewWorkflowRepository()

	transferService := service.NewTransferService(
		db,
		transferRepository,
		gridRowDataRepositoryWithChild,
	)

	uploadService := service.NewUploadService(
		db,
		grTransferRepositoryWithChild,
		userRepository,
		workflowRepository,
		transferRepository,
		inventoryRepository,
		documentRepository,
		AccountID,
		"en",
	)
	handler := &handler.HTTPTreeGridHandler{
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := transferService.GetPageCount(tr)
			return float64(count), err
		},
		CallbackGetPageDataFunc: transferService.GetPageData,
		CallbackUploadDataFunc: func(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
			return uploadService.Handle(req)
		},
	}

	return handler
}
