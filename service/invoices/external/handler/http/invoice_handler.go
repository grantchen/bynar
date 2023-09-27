package http_handler

import (
	"database/sql"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func NewHTTPHandler(appConfig config.AppConfig, db *sql.DB) *handler.HTTPTreeGridHandler {
	simpleInvoicesRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "invoicess", repository.InvoiceFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
	treeGridService, _ := service.NewTreeGridService(db, simpleInvoicesRepository, 0)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:   treeGridService.Handle,
		CallbackGetPageDataFunc:  treeGridService.GetPageData,
		CallbackGetPageCountFunc: treeGridService.GetPageCount,
	}
	return handler
}
