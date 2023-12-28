package service

import (
	"context"
	"database/sql"

	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/service"
)

type treegridService struct {
	db                  *sql.DB
	procurementsService service.ProcurementsService
	uploadService       *service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.Service {
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
	procurementRepository := pkgrepository.NewProcurementRepository(db)
	unitRepo := pkgrepository.NewUnitRepository(db)
	currencyRepo := pkgrepository.NewCurrencyRepository(db)
	invRepo := pkgrepository.NewInventoryRepository(db)
	documentRepository := pkgrepository.NewDocuments(db, "procurements")
	workflowRepository := pkgrepository.NewWorkflowRepository(db)
	// init services
	approvalSvc := pkgservice.NewApprovalCashPaymentService(pkgrepository.NewApprovalOrder(
		workflowRepository,
		procurementRepository),
	)

	docSvc := pkgservice.NewDocumentService(documentRepository)

	procService := service.NewProcurementSvc(db, gridRowDataRepositoryWithChild, procurementRepository, unitRepo, currencyRepo, invRepo, language)

	grPaymentDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames)

	uploadService := service.NewUploadService(db, language, approvalSvc, docSvc, grPaymentDataUploadRepositoryWithChild, accountID, procService)
	return &treegridService{
		db:                  db,
		procurementsService: procService,
		uploadService:       uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (*treegridService) GetCellData(_ context.Context, _ *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.procurementsService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.procurementsService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
