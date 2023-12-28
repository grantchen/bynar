package service

import (
	"context"
	"database/sql"

	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/service"
)

// treegridService implements treegrid.Service
type treegridService struct {
	db            *sql.DB
	saleService   service.SaleService
	uploadService *service.UploadService
}

// newTreeGridService returns new treegrid.Service
func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.Service {
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
	workflowRepository := pkgrepository.NewWorkflowRepository(db)
	unitRepository := pkgrepository.NewUnitRepository(db)
	currencyRepository := pkgrepository.NewCurrencyRepository(db)
	inventoryRepository := pkgrepository.NewInventoryRepository(db)
	boundFlowRepository := pkgrepository.NewBoundFlows()

	documentRepository := pkgrepository.NewDocuments(db, "procurements")

	approvalSvc := pkgservice.NewApprovalCashPaymentService(pkgrepository.NewApprovalOrder(
		workflowRepository,
		saleRepository),
	)
	docSvc := pkgservice.NewDocumentService(documentRepository)

	uploadService := service.NewUploadService(
		db,
		grSaleRepository,
		grSaleDataUploadRepositoryWithChild,
		language,
		accountID,
		approvalSvc,
		docSvc,
		saleRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository,
		boundFlowRepository,
	)
	return &treegridService{
		db:            db,
		saleService:   saleService,
		uploadService: uploadService,
	}
}

// NewTreeGridServiceFactory returns new treegrid.TreeGridServiceFactoryFunc
func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (s *treegridService) GetCellData(_ context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	return s.saleService.GetCellSuggestion(req)
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.saleService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.saleService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
