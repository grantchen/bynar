package main

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

const AccountId = 1

func main() {

	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()
	connString := appConfig.GetAccountManagementConnection()
	connAccountString := appConfig.GetAccountManagementConnection()
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	connString = "2Cef1WJMcYBTnfh.root:mEt9041wub4bdbCW@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/c8460451-c997-4289-ad09-0e4b74ef7fb2?tls=true"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	dbAccount, _ := sql_db.NewConnection(connAccountString)

	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"payments",
		"payment_lines",
		repository.PaymentFieldNames,
		repository.PaymentLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "document_id",
			QueryParent:              repository.QueryParent,
			QueryParentCount:         repository.QueryParentCount,
			QueryParentJoins:         repository.QueryParentJoins,
			QueryChild:               repository.QueryChild,
			QueryChildCount:          repository.QueryChildCount,
			QueryChildJoins:          repository.QueryChildJoins,
			ChildJoinFieldWithParent: "parent_id",
			ParentIdField:            "id",
		},
	)
	paymentRepository := repository.NewPayment(db, "payments", "payment_lines")
	procurementRepository := pkg_repository.NewProcurementRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	cashManagementRepository := pkg_repository.NewCashManagementRepository(db)

	paymentService := service.NewPaymentService(db, gridRowDataRepositoryWithChild, paymentRepository, procurementRepository, currencyRepository, cashManagementRepository)

	grPaymentDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "payments",
		"payment_lines",
		repository.PaymentFieldNames,
		repository.PaymentLineFieldNames)

	grPaymentLineRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"payment_lines",
		repository.PaymentLineFieldNames,
		1, // arbitrary
	)
	grPaymentRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"payments",
		repository.PaymentFieldNames,
		1, // arbitrary
	)
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	documentRepository := pkg_repository.NewDocuments(db, "procurements")
	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
		workflowRepository,
		paymentRepository),
	)

	docSvc := pkg_service.NewDocumentService(documentRepository)
	//todo refactor accountID
	uploadService := service.NewUploadService(db, grPaymentRepository, grPaymentDataUploadRepositoryWithChild, grPaymentLineRepository, "en", approvalSvc, docSvc, AccountId, paymentService)

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: paymentService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := paymentService.GetPageCount(tr)
			return float64(count), err
		},
		AccountManagerService: accountService,
	}

	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)
	http.HandleFunc("/cell", handler.HTTPHandleCell)

	// server
	log.Println("start server at 8081!")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
