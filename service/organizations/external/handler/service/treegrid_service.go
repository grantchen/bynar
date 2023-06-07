package service

import (
	"context"
	"database/sql"
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db                  *sql.DB
	organizationService service.OrganizationService
	uploadService       service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int) treegrid.TreeGridService {

	var filterPermissionCondition string

	logger.Debug("accountID:", accountID)
	if accountID != 0 {
		filterPermissionCondition = fmt.Sprintf(repository.QueryPermissionFormat, accountID)
	}

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "organizations", repository.OrganizationFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "code",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: filterPermissionCondition,
		})
	organizationService := service.NewOrganizationService(db, simpleOrganizationRepository)

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository)
	return &treegridService{
		db:                  db,
		organizationService: organizationService,
		uploadService:       *uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, permissionInfo *treegrid.PermissionInfo) treegrid.TreeGridService {
		return newTreeGridService(db, accountID)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) float64 {
	return float64(s.organizationService.GetPageCount(tr))
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.organizationService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
