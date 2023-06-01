package main

import (
	"fmt"
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/external/service"
)

func main() {
	secretmanager, err := utils.GetSecretManager()
	if err != nil {
		fmt.Printf("error: %v", err)
		log.Panic(err)
	}

	connAccountString := "root:123456@tcp(localhost:3306)/accounts_management"
	dbAccount, err := sql_db.NewConnection(connAccountString)

	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, secretmanager)

	treegridService := service.NewTreeGridServiceFactory()
	handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountService,
		TreeGridServiceFactory: treegridService,
	}

	handler.HandleHTTPReqWithAuthenMWAndDefaultPath()

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
