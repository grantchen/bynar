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
	ModuleID  int = 4
	AccountID int = 123456
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
	gridRowRep := treegrid.NewGridRepository(db, "procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames)

	procurementRepository := pkg_repository.NewProcurementRepository(db)
	workflowRepository := pkg_repository.NewWorkflowRepository(db, ModuleID)
	unitRepository := pkg_repository.NewUnitRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	inventoryRepository := pkg_repository.NewInventoryRepository(db)

	documentRepository := pkg_repository.NewDocuments(db, "procurements")

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		procurementRepository),
	)

	procrSvc := service.NewProcurementSvc(
		procurementRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository)

	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadSvc, _ := service.NewService(
		db,
		approvalSvc,
		gridRowRep,
		procrSvc,
		ModuleID,
		AccountID,
		docSvc,
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: uploadSvc.Handle,
	}

	return handler
}
