package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// treegridService implements treegrid.Service
type treegridService struct {
	db            *sql.DB
	siteService   service.SiteService
	uploadService service.UploadService
}

// newTreeGridService create new treegridService
func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.Service {
	logger.Debug("accountID:", accountID)
	simpleSiteRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "sites", repository.SiteFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "code",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: repository.QueryPermissionFormat,
		})
	siteService := service.NewSiteService(db, simpleSiteRepository)

	uploadService, _ := service.NewUploadService(db, siteService, simpleSiteRepository, language)
	return &treegridService{
		db:            db,
		siteService:   siteService,
		uploadService: *uploadService,
	}
}

// NewTreeGridServiceFactory create new treegrid.TreeGridServiceFactoryFunc
func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, siteUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (*treegridService) GetCellData(_ context.Context, _ *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.siteService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.siteService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
