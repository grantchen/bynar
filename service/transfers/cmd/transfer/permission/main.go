package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/external/handler/service"
)

func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()

	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	dbAccount, err := sql_db.NewConnection(accountManagementConnectionString)

	if err != nil {
		log.Panic(err)
	}

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	connectionPool := connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	treegridService := service.NewTreeGridServiceFactory()
	handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountService,
		TreeGridServiceFactory: treegridService,
		ConnectionPool:         connectionPool,
		PathPrefix:             "/apprunnerurl/organizations",
	}

	handler.HandleHTTPReqWithAuthenMWAndDefaultPath()

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}