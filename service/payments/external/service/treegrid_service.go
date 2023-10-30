package service

import (
	"context"
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db             *sql.DB
	paymentService service.PaymentService
	uploadService  *service.UploadService
	accountId      int
}

func newTreeGridService(db *sql.DB, language string, accountId int) treegrid.TreeGridService {
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
			QueryChildSuggestion:     repository.QueryChildSuggestion,
			ChildJoinFieldWithParent: "parent_id",
			ParentIdField:            "id",
		},
	)

	paymentRepository := repository.NewPayment(db, "payments", "payment_lines")
	procurementRepository := pkg_repository.NewProcurementRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	cashManagementRepository := pkg_repository.NewCashManagementRepository(db)
	documentRepository := pkg_repository.NewDocuments(db, "procurements")
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	// init services
	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		paymentRepository),
	)

	docSvc := pkg_service.NewDocumentService(documentRepository)

	paymentService := service.NewPaymentService(db, gridRowDataRepositoryWithChild, paymentRepository, procurementRepository, currencyRepository, cashManagementRepository)

	grPaymentDataUploadRepositoryWithChild := treegrid.NewGridRepository(db,
		"payments",
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

	uploadService := service.NewUploadService(db, grPaymentRepository, grPaymentDataUploadRepositoryWithChild, grPaymentLineRepository, language, approvalSvc, docSvc, accountId, paymentService)
	return &treegridService{
		db:             db,
		paymentService: paymentService,
		uploadService:  uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, AccountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, language, AccountID)
	}
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.paymentService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.paymentService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}

// GetCellData implements treegrid.TreeGridService
func (s *treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}
