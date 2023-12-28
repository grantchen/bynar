package service

import (
	"context"
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/service"
)

type treegridService struct {
	db                *sql.DB
	warehousesService service.WarehousesService
	uploadService     service.UploadService
}

func newTreeGridService(db *sql.DB, _ int, language string) treegrid.Service {
	simpleWarehousesRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(
		db,
		"warehouses",
		repository.WarehousesFieldNames,
		100,
		&treegrid.SimpleGridRepositoryCfg{
			MainCol:     "code",
			QueryString: repository.QuerySelect,
			QueryJoin:   repository.QueryJoin,
			QueryCount:  repository.QueryCount,
		},
	)

	warehousesRepository := repository.NewPostingSetupRepository(db)
	warehousesService := service.NewWarehousesService(simpleWarehousesRepository)
	uploadService := service.NewUploadService(
		db,
		simpleWarehousesRepository,
		warehousesRepository,
		language,
	)
	return &treegridService{
		db:                db,
		warehousesService: warehousesService,
		uploadService:     uploadService,
	}
}

func NewTreeGridServiceFactory() treegrid.ServiceFactoryFunc {
	return func(db *sql.DB, accountID int, organizationUuid string, permissionInfo *treegrid.PermissionInfo, language string) treegrid.Service {
		return newTreeGridService(db, accountID, language)
	}
}

// GetCellData implements treegrid.Service
func (*treegridService) GetCellData(context.Context, *treegrid.Treegrid) (*treegrid.PostResponse, error) {
	panic("unimplemented")
}

// GetPageCount implements treegrid.Service
func (s *treegridService) GetPageCount(tr *treegrid.Treegrid) (float64, error) {
	count, err := s.warehousesService.GetPageCount(tr)
	return float64(count), err
}

// GetPageData implements treegrid.Service
func (s *treegridService) GetPageData(tr *treegrid.Treegrid) ([]map[string]string, error) {
	return s.warehousesService.GetPageData(tr)
}

// Upload implements treegrid.Service
func (s *treegridService) Upload(req *treegrid.PostRequest) (*treegrid.PostResponse, error) {
	return s.uploadService.Handle(req)
}
