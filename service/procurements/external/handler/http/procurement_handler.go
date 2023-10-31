package http_handler

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/service"
)

// TODO: get throug request
var (
	AccountID int = 123456
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "document_id",
			QueryParent:              repository.QueryParent,
			QueryParentCount:         repository.QueryParentCount,
			QueryParentJoins:         repository.QueryParentJoins,
			QueryChild:               repository.QueryChild,
			QueryChildCount:          repository.QueryChildCount,
			QueryChildJoins:          repository.QueryChildJoins,
			ChildJoinFieldWithParent: "parent_id",
			ParentIdField:            "id",
		},
	)
	gridRowRep := treegrid.NewGridRepository(db, "procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames)

	procurementRepository := pkg_repository.NewProcurementRepository(db)
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	unitRepository := pkg_repository.NewUnitRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	inventoryRepository := pkg_repository.NewInventoryRepository(db)

	documentRepository := pkg_repository.NewDocuments(db, "procurements")

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		procurementRepository),
	)

	procrSvc := service.NewProcurementSvc(
		db,
		gridRowDataRepositoryWithChild,
		procurementRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository)

	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadSvc := service.NewUploadService(
		db,
		"en",
		approvalSvc,
		docSvc,
		gridRowRep,
		AccountID,
		procrSvc,
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadSvc.Handle,
		CallbackGetPageDataFunc: procrSvc.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := procrSvc.GetPageCount(tr)
			return float64(count), err
		},
	}

	return handler
}
