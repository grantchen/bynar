package factory

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/mapping"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repositories"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repositories/gridrep"
	svc "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"
	svc_upload "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service/upload"
)

func NewUploadService(sqlDB *sql.DB, moduleID int, accountID int) (svc_upload.UploadService, error) {
	// init repositories
	// grid repo
	gridRowRepository := gridrep.NewGridRepository(
		sqlDB,
		"payments",
		"payment_lines",
		mapping.PaymentFieldNames,
		mapping.PaymentLineFieldNames,
	)

	paymentRepository := repositories.NewPayment(sqlDB, "payments", "payment_lines")
	procurementRepository := repositories.NewProcurementRepository(sqlDB)
	currencyRepository := repositories.NewCurrencyRepository(sqlDB)
	cashManagementRepository := repositories.NewCashManagementRepository(sqlDB)
	documentRepository := repositories.NewDocuments(sqlDB, "procurements")
	workflowRepository := repositories.NewWorkflowRepository(sqlDB, moduleID)
	// init services
	// TODO

	approvalSvc := svc.NewApprovalCashPaymentService(repositories.NewApprovalOrder(
		workflowRepository,
		paymentRepository),
	)

	paymentSvc := svc.NewPaymentService(
		paymentRepository,
		procurementRepository,
		currencyRepository,
		cashManagementRepository,
	)

	docSvc := svc.NewDocumentService(documentRepository)

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
