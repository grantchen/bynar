package http_handler

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

const AccountId = 1

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {

	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"payments",
		"payment_lines",
		repository.PaymentFieldNames,
		repository.PaymentLineFieldNames,
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
	paymentRepository := repository.NewPayment(db, "payments", "payment_lines")
	procurementRepository := pkg_repository.NewProcurementRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	cashManagementRepository := pkg_repository.NewCashManagementRepository(db)
	// todo refactor moduleId
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	documentRepository := pkg_repository.NewDocuments(db, "procurements")
	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		paymentRepository),
	)

	docSvc := pkg_service.NewDocumentService(documentRepository)
	paymentService := service.NewPaymentService(db, gridRowDataRepositoryWithChild, paymentRepository, procurementRepository, currencyRepository, cashManagementRepository)

	grPaymentDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "payments",
		"payment_lines",
		repository.PaymentFieldNames,
		repository.PaymentLineFieldNames)

	grPaymentLineRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"payment_lines",
		repository.PaymentLineFieldNames,
		1, // arbitrary
	)
	grPaymentRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"payments",
		repository.PaymentFieldNames,
		1, // arbitrary
	)

	uploadService := service.NewUploadService(db, grPaymentRepository, grPaymentDataUploadRepositoryWithChild, grPaymentLineRepository, "en", approvalSvc, docSvc, AccountId, paymentService)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: paymentService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := paymentService.GetPageCount(tr)
			return float64(count), err
		},
	}

	return handler
}
