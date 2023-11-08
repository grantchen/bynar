package http_handler

import (
	"database/sql"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/service"
)

// for test
const accountId = 6

// NewHTTPHandler returns new http handler
func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"sales",
		"sale_lines",
		repository.SaleFieldNames,
		repository.SaleLineFieldNames,
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
	saleService := service.NewSaleService(db, gridRowDataRepositoryWithChild)

	grSaleDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "sales",
		"sale_lines",
		repository.SaleFieldNames,
		repository.SaleLineFieldNames)

	grSaleRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"sales",
		repository.SaleFieldNames,
		1, // arbitrary
	)

	saleRepository := repository.NewSaleRepository(db)
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	unitRepository := pkg_repository.NewUnitRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	inventoryRepository := pkg_repository.NewInventoryRepository(db)
	boundFlowRepository := pkg_repository.NewBoundFlows()

	documentRepository := pkg_repository.NewDocuments(db, "procurements")

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		saleRepository),
	)
	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadService := service.NewUploadService(
		db,
		grSaleRepository,
		grSaleDataUploadRepositoryWithChild,
		"en",
		accountId,
		approvalSvc,
		docSvc,
		saleRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository,
		boundFlowRepository,
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: saleService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := saleService.GetPageCount(tr)
			return float64(count), err
		},
		CallBackGetCellDataFunc: saleService.GetCellSuggestion,
	}

	return handler
}
