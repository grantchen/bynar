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

type treegridService struct {
	db            *sql.DB
	uploadService service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int) treegrid.TreeGridService {
	logger.Debug("accountID:", accountID)

	var additionWhere string
	if accountID > 0 {
		additionWhere = fmt.Sprintf(repository.AdditionWhere, accountID)
	}

	simpleInvoiceRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "invoices", repository.InvoiceFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "code",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: additionWhere,
		})
	uploadService, _ := service.NewUploadService(db, simpleInvoiceRepository, accountID)
	return &treegridService{
		db:            db,
		uploadService: *uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo) treegrid.TreeGridService {
		return newTreeGridService(db, accountID)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.uploadService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.uploadService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
