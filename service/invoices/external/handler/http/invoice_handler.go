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

	uploadService, _ := service.NewUploadService(db, simpleInvoicesRepository, 0)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: uploadService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := uploadService.GetPageCount(tr)
			return float64(count), err
		},
	}
	return handler

}
