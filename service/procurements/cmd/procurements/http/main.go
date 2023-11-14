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
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/internal/service"
)

// TODO: get throug request
var (
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

	procurementRepository := pkg_repository.NewProcurementRepository(db)
	workflowRepository := pkg_repository.NewWorkflowRepository(db)
	unitRepository := pkg_repository.NewUnitRepository(db)
	currencyRepository := pkg_repository.NewCurrencyRepository(db)
	inventoryRepository := pkg_repository.NewInventoryRepository(db)

	documentRepository := pkg_repository.NewDocuments(db, "procurements")

	approvalSvc := pkg_service.NewApprovalCashPaymentService(pkg_repository.NewApprovalOrder(
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

	docSvc := pkg_service.NewDocumentService(documentRepository)

	uploadSvc := service.NewUploadService(
		db,
		"en",
		approvalSvc,
		docSvc,
		gridRowRep,
		AccountID,
		procrSvc,
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: uploadSvc.Handle,
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
