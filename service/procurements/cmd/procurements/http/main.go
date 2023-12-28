package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/service"
)

// AccountID TODO: get throug request
var (
	AccountID = 123456
)

func main() {
	secretmanager, err := utils.GetSecretManager()
	if err != nil {
		log.Panic(err)
	}

	appConfig := config.NewAWSSecretsManagerConfig(secretmanager)
	connString := appConfig.GetDBConnection()
	db, err := sqldb.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames,
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

	gridRowRep := treegrid.NewGridRepository(db, "procurements",
		"procurement_lines",
		repository.ProcurementFieldNames,
		repository.ProcurementLineFieldNames)

	procurementRepository := pkgrepository.NewProcurementRepository(db)
	workflowRepository := pkgrepository.NewWorkflowRepository(db)
	unitRepository := pkgrepository.NewUnitRepository(db)
	currencyRepository := pkgrepository.NewCurrencyRepository(db)
	inventoryRepository := pkgrepository.NewInventoryRepository(db)

	documentRepository := pkgrepository.NewDocuments(db, "procurements")

	approvalSvc := pkgservice.NewApprovalCashPaymentService(pkgrepository.NewApprovalOrder(
		workflowRepository,
		procurementRepository),
	)

	procrSvc := service.NewProcurementSvc(
		db,
		gridRowDataRepositoryWithChild,
		procurementRepository,
		unitRepository,
		currencyRepository,
		inventoryRepository,
		"en")

	docSvc := pkgservice.NewDocumentService(documentRepository)

	uploadSvc := service.NewUploadService(
		db,
		"en",
		approvalSvc,
		docSvc,
		gridRowRep,
		AccountID,
		procrSvc,
	)

	h := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: uploadSvc.Handle,
	}
	http.HandleFunc("/upload", h.HTTPHandleUpload)
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
