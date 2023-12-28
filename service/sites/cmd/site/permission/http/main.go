package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/external/handler/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
)

func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()

	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	dbAccount, err := sqldb.NewConnection(accountManagementConnectionString)

	if err != nil {
		log.Panic(err)
	}

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkgrepository.NewAccountManagerRepository(dbAccount)
	accountService := pkgservice.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	connectionPool := connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	treegridService := service.NewTreeGridServiceFactory()
	h := &handler.HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountService,
		TreeGridServiceFactory: treegridService,
		ConnectionPool:         connectionPool,
		PathPrefix:             "/apprunnerurl/sites",
	}

	h.HandleHTTPReqWithAuthenMWAndDefaultPath()

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
