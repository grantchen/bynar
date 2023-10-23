package http_handler

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
	simpleSiteRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "sites", repository.SiteFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
	siteService := service.NewSiteService(db, simpleSiteRepository)

	uploadService, _ := service.NewUploadService(db, siteService, simpleSiteRepository, "en")

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: siteService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := siteService.GetPageCount(tr)
			return float64(count), err
		},
	}
	return handler

}
