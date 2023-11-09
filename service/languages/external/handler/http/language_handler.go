package http_handler

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
	simpleLanguageRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "languages", repository.LanguageFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
	languageService := service.NewLanguageService(db, simpleLanguageRepository)

	uploadService, _ := service.NewUploadService(db, languageService, simpleLanguageRepository, "en")

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: languageService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := languageService.GetPageCount(tr)
			return float64(count), err
		},
	}
	return handler

}
