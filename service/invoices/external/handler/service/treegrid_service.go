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

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.Service {
	logger.Debug("accountID:", accountID)

	simpleInvoiceRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "invoices", repository.InvoiceFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "invoice_date",
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

func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (*treeGridService) GetCellData(_ context.Context, _ *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treeGridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	return s.internalTreeGridService.GetPageCount(tr)
}

// GetPageData implements treegrid.Service
func (s *treeGridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.internalTreeGridService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treeGridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.internalTreeGridService.Handle(req)
}
