package http_handler

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/service"
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
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
		"en",
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: warehousesService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := warehousesService.GetPageCount(tr)
			return float64(count), err
		},
	}
	//http.HandleFunc("/upload", handler.HTTPHandleUpload)
	//http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	//http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	return handler
}
