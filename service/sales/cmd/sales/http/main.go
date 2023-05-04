package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/service"
)

// TODO: get throug request
var (
	ModuleID  int = 4
	AccountID int = 123456
)

func main() {
	secretmanager, err := utils.GetSecretManager()
	if err != nil {
		log.Panic(err)
	}

	appConfig := config.NewAWSSecretsManagerConfig(secretmanager)
	connString := appConfig.GetDBConnection()
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	gridRowRep := treegrid.NewGridRepository(db, "sales",
		"sale_lines",
		repository.SaleFieldNames,
		repository.SaleLineFieldNames)

	saleRepository := repository.NewSaleRepository(db)
	workflowRepository := pkg_repository.NewWorkflowRepository(db, ModuleID)
	unitRepository := pkg_repository.NewUnitRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	inventoryRepository := pkg_repository.NewInventoryRepository(db)

	documentRepository := pkg_repository.NewDocuments(db, "procurements")

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		saleRepository),
	)

	saleService := service.NewSaleService(
		saleRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository)

	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadSvc, err := service.NewService(
		db,
		approvalSvc,
		gridRowRep,
		saleService,
		ModuleID,
		AccountID,
		docSvc,
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: uploadSvc.Handle,
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
