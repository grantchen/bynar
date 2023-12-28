package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"

	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
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
	db, err := sqldb.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	dbAccount, _ := sqldb.NewConnection(connAccountString)

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
	workflowRepository := pkgrepository.NewWorkflowRepository(db)
	unitRepository := pkgrepository.NewUnitRepository(db)
	currencyRepository := pkgrepository.NewCurrencyRepository(db)
	inventoryRepository := pkgrepository.NewInventoryRepository(db)
	boundFlowRepository := pkgrepository.NewBoundFlows()

	documentRepository := pkgrepository.NewDocuments(db, "procurements")

	approvalSvc := pkgservice.NewApprovalCashPaymentService(pkgrepository.NewApprovalOrder(
		workflowRepository,
		saleRepository),
	)
	docSvc := pkgservice.NewDocumentService(documentRepository)

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

	accountRepository := pkgrepository.NewAccountManagerRepository(dbAccount)
	accountService := pkgservice.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	h := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: saleService.GetPageData,
		CallBackGetCellDataFunc: saleService.GetCellSuggestion,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := saleService.GetPageCount(tr)
			return float64(count), err
		},
		AccountManagerService: accountService,
	}

	http.HandleFunc("/upload", h.HTTPHandleUpload)
	http.HandleFunc("/data", h.HTTPHandleGetPageCount)
	http.HandleFunc("/page", h.HTTPHandleGetPageData)
	http.HandleFunc("/cell", h.HTTPHandleCell)

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
