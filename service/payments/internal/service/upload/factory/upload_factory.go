package factory

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/mapping"
	repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repositories"
	svc "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"
	svc_upload "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewUploadService(sqlDB *sql.DB, moduleID int, accountID int) (svc_upload.UploadService, error) {
	// init repositories
	// grid repo
	gridRowRepository := treegrid.NewGridRepository(
		sqlDB,
		"payments",
		"payment_lines",
		mapping.PaymentFieldNames,
		mapping.PaymentLineFieldNames,
	)

	paymentRepository := repository.NewPayment(sqlDB, "payments", "payment_lines")
	procurementRepository := pkg_repository.NewProcurementRepository(sqlDB)
	currencyRepository := pkg_repository.NewCurrencyRepository(sqlDB)
	cashManagementRepository := pkg_repository.NewCashManagementRepository(sqlDB)
	documentRepository := pkg_repository.NewDocuments(sqlDB, "procurements")
	workflowRepository := pkg_repository.NewWorkflowRepository(sqlDB, moduleID)
	// init services
	// TODO

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		paymentRepository),
	)

	paymentSvc := svc.NewPaymentService(
		paymentRepository,
		procurementRepository,
		currencyRepository,
		cashManagementRepository,
	)

	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadSvc, err := svc_upload.NewUploadService(
		sqlDB,
		approvalSvc,
		gridRowRepository,
		paymentSvc,
		moduleID,
		accountID,
		docSvc)

	return uploadSvc, err
}
