package service

import (
	"context"
	"database/sql"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treeGridService struct {
	db                      *sql.DB
	internalTreeGridService service.TreeGridService
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.TreeGridService {
	logger.Debug("accountID:", accountID)

	simpleInvoiceRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "invoices", repository.InvoiceFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "code",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: fmt.Sprintf(repository.AdditionWhere, accountID),
		})
	internalTreeGridService, _ := service.NewTreeGridService(db, simpleInvoiceRepository, accountID, language)

	return &treeGridService{
		db:                      db,
		internalTreeGridService: *internalTreeGridService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treeGridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treeGridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	return s.internalTreeGridService.GetPageCount(tr)
}

// GetPageData implements treegrid.TreeGridService
func (s *treeGridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.internalTreeGridService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treeGridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.internalTreeGridService.Handle(req)
}
