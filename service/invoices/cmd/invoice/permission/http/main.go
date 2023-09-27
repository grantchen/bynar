package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/external/handler/service"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
)

func main() {
	connAccountString := "root:123456@tcp(localhost:3306)/accounts_management"
	dbAccount, err := sql_db.NewConnection(connAccountString)
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

	handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountService,
		TreeGridServiceFactory: service.NewTreeGridServiceFactory(),
		ConnectionPool:         connectionPool,
		PathPrefix:             "/apprunnerurl/invoices",
		IsValidatePermissions:  false,
	}
	handler.HandleHTTPReqWithAuthenMWAndDefaultPath()

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
