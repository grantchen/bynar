package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db              *sql.DB
	transferService service.TransferService
	accountID       int
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.TreeGridService {
	logger.Debug("accountID:", accountID)

	documentRepository := repository.NewDocumentRepository(db)
	inventoryRepository := repository.NewInventoryRepository(db)
	transferRepository := repository.NewTransferRepository(db)
	userRepository := repository.NewUserRepository(db)
	workflowRepository := repository.NewWorkflowRepository()

	transferService := service.NewTransferService(
		db,
		language,
		userRepository,
		workflowRepository,
		transferRepository,
		inventoryRepository,
		documentRepository,
	)

	return &treegridService{
		db:              db,
		accountID:       accountID,
		transferService: transferService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, transferUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.transferService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.transferService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.transferService.HandleUpload(req, s.accountID)
}
