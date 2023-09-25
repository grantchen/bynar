package http_handler

import (
	"database/sql"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "organizations", repository.OrganizationFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
	organizationService := service.NewOrganizationService(db, simpleOrganizationRepository)

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: organizationService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := organizationService.GetPageCount(tr)
			return float64(count), err
		},
	}
	return handler

}
