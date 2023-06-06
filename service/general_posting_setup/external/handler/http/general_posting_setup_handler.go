package http_handler

import (
	"database/sql"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
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
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: generalPostingSetupService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) float64 {
			return float64(generalPostingSetupService.GetPageCount(tr))
		},
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	return handler
}
