package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"

	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
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
	connString = "2Cef1WJMcYBTnfh.root:mEt9041wub4bdbCW@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/bynar_dev?tls=true"
	db, err := sqldb.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	dbAccount, _ := sqldb.NewConnection(connAccountString)

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
	procurementRepository := pkgrepository.NewProcurementRepository(db)
	currencyRepository := pkgrepository.NewCurrencyRepository(db)
	cashManagementRepository := pkgrepository.NewCashManagementRepository(db)

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
	workflowRepository := pkgrepository.NewWorkflowRepository(db)
	documentRepository := pkgrepository.NewDocuments(db, "procurements")
	approvalSvc := pkgservice.NewApprovalCashPaymentService(pkgrepository.NewApprovalOrder(
		workflowRepository,
		paymentRepository),
	)

	docSvc := pkgservice.NewDocumentService(documentRepository)
	//todo refactor accountID
	uploadService := service.NewUploadService(db, grPaymentRepository, grPaymentDataUploadRepositoryWithChild, grPaymentLineRepository, "en", approvalSvc, docSvc, AccountId, paymentService)

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkgrepository.NewAccountManagerRepository(dbAccount)
	accountService := pkgservice.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	h := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: paymentService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := paymentService.GetPageCount(tr)
			return float64(count), err
		},
		AccountManagerService: accountService,
	}

	http.HandleFunc("/upload", h.HTTPHandleUpload)
	http.HandleFunc("/data", h.HTTPHandleGetPageCount)
	http.HandleFunc("/page", h.HTTPHandleGetPageData)
	http.HandleFunc("/cell", h.HTTPHandleCell)

	// server
	log.Println("start server at 8081!")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
