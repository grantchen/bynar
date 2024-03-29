package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// treegridService implements treegrid.Service
type treegridService struct {
	db              *sql.DB
	transferService service.TransferService
	accountID       int
	uploadService   *service.UploadService
}

// newTreeGridService returns a new treegridService
func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.Service {
	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"transfers",
		"transfer_lines",
		repository.TransferFieldNames,
		repository.TransferLineFieldNames,
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

	grTransferRepositoryWithChild := treegrid.NewGridRepository(db,
		"transfers",
		"transfer_lines",
		repository.TransferFieldNames,
		repository.TransferLineFieldNames,
	)

	documentRepository := repository.NewDocumentRepository(db)
	inventoryRepository := repository.NewInventoryRepository(db)
	transferRepository := repository.NewTransferRepository(db, language)
	userRepository := repository.NewUserRepository(db)

	transferService := service.NewTransferService(
		db,
		transferRepository,
		gridRowDataRepositoryWithChild,
	)

	uploadService := service.NewUploadService(
		db,
		grTransferRepositoryWithChild,
		userRepository,
		transferRepository,
		inventoryRepository,
		documentRepository,
		accountID,
		language,
	)
	return &treegridService{
		db:              db,
		transferService: transferService,
		uploadService:   uploadService,
	}
}

// NewTreeGridServiceFactory returns a new treegrid.TreeGridServiceFactoryFunc
func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, transferUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (*treegridService) GetCellData(context.Context, *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.transferService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.transferService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
