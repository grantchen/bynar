package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/internal/service"
)

// for test
const accountId = 6

func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()
	connString := appConfig.GetAccountManagementConnection()
	connAccountString := appConfig.GetAccountManagementConnection()
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	connString = "root:123456@tcp(localhost:3306)/46542255-9d45-49d5-939d-84bc55b1a938"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	dbAccount, _ := sql_db.NewConnection(connAccountString)

	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"sales",
		"sale_lines",
		repository.SaleFieldNames,
		repository.SaleLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "document_id",
			QueryParent:              repository.QueryParent,
			QueryParentCount:         repository.QueryParentCount,
			QueryParentJoins:         repository.QueryParentJoins,
			QueryChild:               repository.QueryChild,
			QueryChildCount:          repository.QueryChildCount,
			QueryChildJoins:          repository.QueryChildJoins,
			QueryChildSuggestion:     repository.QueryChildSuggestion,
			ChildJoinFieldWithParent: "parent_id",
			ParentIdField:            "id",
		},
	)
	saleService := service.NewSaleService(db, gridRowDataRepositoryWithChild)

	grSaleDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "sales",
		"sale_lines",
		repository.SaleFieldNames,
		repository.SaleLineFieldNames)

	grSaleRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"sales",
		repository.SaleFieldNames,
		1, // arbitrary
	)

	saleRepository := repository.NewSaleRepository(db)
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	unitRepository := pkg_repository.NewUnitRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	inventoryRepository := pkg_repository.NewInventoryRepository(db)
	boundFlowRepository := pkg_repository.NewBoundFlows()

	documentRepository := pkg_repository.NewDocuments(db, "procurements")

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		saleRepository),
	)
	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadService := service.NewUploadService(
		db,
		grSaleRepository,
		grSaleDataUploadRepositoryWithChild,
		"en",
		accountId,
		approvalSvc,
		docSvc,
		saleRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository,
		boundFlowRepository,
	)

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: saleService.GetPageData,
		CallBackGetCellDataFunc: saleService.GetCellSuggestion,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := saleService.GetPageCount(tr)
			return float64(count), err
		},
		AccountManagerService: accountService,
	}

	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)
	http.HandleFunc("/cell", handler.HTTPHandleCell)

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
