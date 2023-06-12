package main

import (
	"fmt"
	"log"
	"net/http"

	general_posting_setup_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/external/handler/http"
	organizations_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/external/handler/service"
	payments_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/external/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	procurements_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/external/handler/http"
	sales_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/external/handler/http"
	transfers_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/external/handler/http"
	usergroups_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/external/handler/http"
)

type HandlerMapping struct {
	handler    *handler.HTTPTreeGridHandler
	prefixPath string
}

type HandlerMappingWithPermission struct {
	factoryFunc treegrid.TreeGridServiceFactoryFunc
	prefixPath  string
}

const prefix = "/apprunnerurl"

func main() {

	secretmanager, err := utils.GetSecretManager()
	if err != nil {
		fmt.Printf("error: %v", err)
		log.Panic(err)
	}

	appConfig := config.NewAWSSecretsManagerConfig(secretmanager)

	connString := appConfig.GetDBConnection()
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	lsHandlerMapping := make([]*HandlerMapping, 0)
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: sales_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/sales"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: payments_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/payments"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: transfers_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/transfers"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: procurements_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/procurements"})
	// lsHandlerMapping = append(lsHandlerMapping,
	// 	&HandlerMapping{handler: organizations_handler.NewHTTPHandler(appConfig, db),
	// 		prefixPath: "/organizations"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: usergroups_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/user_groups"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: general_posting_setup_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/general_posting_setup"})

	for _, handlerMapping := range lsHandlerMapping {
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/data", handlerMapping.handler.HTTPHandleGetPageCount)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/page", handlerMapping.handler.HTTPHandleGetPageData)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/upload", handlerMapping.handler.HTTPHandleUpload)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/cell", handlerMapping.handler.HTTPHandleCell)

	}

	// _______________________________________ COMPONENT WITH PERMISSION______________________________________________
	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	logger.Debug("connection string account: ", accountManagementConnectionString)

	connectionPool := connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()
	dbAccount, err := connectionPool.Get(accountManagementConnectionString)

	if err != nil {
		log.Panic(err)
	}
	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, secretmanager)

	lsHandlerMappingWithPermission := make([]*HandlerMappingWithPermission, 0)
	lsHandlerMappingWithPermission = append(lsHandlerMappingWithPermission,
		&HandlerMappingWithPermission{factoryFunc: organizations_service.NewTreeGridServiceFactory(), prefixPath: "/organizations"})

	for _, handlerMappingWithPermission := range lsHandlerMappingWithPermission {
		handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
			AccountManagerService:  accountService,
			TreeGridServiceFactory: handlerMappingWithPermission.factoryFunc,
			ConnectionPool:         connectionPool,
			PathPrefix:             prefix + "/organizations",
		}
		handler.HandleHTTPReqWithAuthenMWAndDefaultPath()
	}

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
