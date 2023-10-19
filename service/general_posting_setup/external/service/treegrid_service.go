package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

type treegridService struct {
	db                         *sql.DB
	generalPostingSetupService service.GeneralPostingSetupService
	uploadService              service.UploadService
}

func newTreeGridService(db *sql.DB, accountID int, language string) treegrid.TreeGridService {
	simpleGeneralPostingSetupRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(
		db,
		"general_posting_setup",
		repository.GeneralPostingSetupFieldNames,
		100,
		&treegrid.SimpleGridRepositoryCfg{
			MainCol:     "code",
			QueryString: repository.QuerySelect,
			QueryJoin:   repository.QueryJoin,
			QueryCount:  repository.QueryCount,
		},
	)

	generalPostingSetupRepository := repository.NewPostingSetupRepository(db)
	generalPostingSetupService := service.NewGeneralPostingSetupService(simpleGeneralPostingSetupRepository)
	uploadService := service.NewUploadService(
		db,
		simpleGeneralPostingSetupRepository,
		generalPostingSetupRepository,
		language,
	)
	return &treegridService{
		db:                         db,
		generalPostingSetupService: generalPostingSetupService,
		uploadService:              uploadService,
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
	count, err := s.generalPostingSetupService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.TreeGridService
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.generalPostingSetupService.GetPageData(tr)
}

// Upload implements treegrid.TreeGridService
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
