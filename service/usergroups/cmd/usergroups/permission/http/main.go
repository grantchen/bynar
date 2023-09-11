package main

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/external/service"
)

func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()

	dbAccount, err := sql_db.NewConnection(appConfig.GetAccountManagementConnection())

	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, nil)

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
		PathPrefix:             "/user_groups",
	}

	handler.HandleHTTPReqWithAuthenMWAndDefaultPath()

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
