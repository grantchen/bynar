package http_handler

import (
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewHTTPHandler(appConfig config.AppConfig) *handler.HTTPTreeGridHandler {
	connString := appConfig.GetDBConnection()
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepository(db, "organizations", repository.OrganizationFieldNames, 100)
	organizationService := service.NewOrganizationService(db, simpleOrganizationRepository)

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: organizationService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) float64 {
			return float64(organizationService.GetPageCount(tr))
		},
	}
	return handler

}
