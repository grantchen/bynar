package service

import (
	"context"
	"database/sql"

	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/service"
)

type treegridService struct {
	db                  *sql.DB
	procurementsService service.ProcurementsService
	uploadService       *service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.TreeGridService {
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
	procurementRepository := pkg_repository.NewProcurementRepository(db)
	unitRepo := pkg_repository.NewUnitRepository(db)
	currencyRepo := pkg_repository.NewCurrencyRepository(db)
	invRepo := pkg_repository.NewInventoryRepository(db)
	documentRepository := pkg_repository.NewDocuments(db, "procurements")
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	// init services
	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		procurementRepository),
	)

	docSvc := pkg_service.NewDocumentService(documentRepository)

	procService := service.NewProcurementSvc(db, gridRowDataRepositoryWithChild, procurementRepository, unitRepo, currencyRepo, invRepo)

	grPaymentDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames)

	uploadService := service.NewUploadService(db, "en", approvalSvc, docSvc, grPaymentDataUploadRepositoryWithChild, accountID, procService)
	return &treegridService{
		db:                  db,
		procurementsService: procService,
		uploadService:       uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.procurementsService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.procurementsService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
