package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db              *sql.DB
	languageService service.LanguageService
	uploadService   service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.TreeGridService {
	logger.Debug("accountID:", accountID)
	simpleLanguageRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "languages", repository.LanguageFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "language",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: repository.QueryPermissionFormat,
		})
	languageService := service.NewLanguageService(db, simpleLanguageRepository)

	uploadService, _ := service.NewUploadService(db, languageService, simpleLanguageRepository, language)
	return &treegridService{
		db:              db,
		languageService: languageService,
		uploadService:   *uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.TreeGridServiceFactoryFunc {
	return func(db *sql.DB, accountID int, languageUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.TreeGridService {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.TreeGridService
func (*treegridService) GetCellData(ctx context.Context, req *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.TreeGridService
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.languageService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.languageService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
