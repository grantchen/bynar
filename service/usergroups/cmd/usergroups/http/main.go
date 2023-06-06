package main

import (
	"fmt"
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/service"
)

func main() {
	// secretmanager, err := utils.GetSecretManager()

	// if err != nil {
	// 	log.Panic(err)
	// }

	// appConfig := config.NewAWSSecretsManagerConfig(secretmanager)
	// connString := appConfig.GetDBConnection()
	// db, err := sql_db.NewConnection(connString)

	// if err != nil {
	// 	log.Panic(err)
	// }

	connString := "root:123456@tcp(localhost:3306)/bynar"
	connAccountString := "root:123456@tcp(localhost:3306)/accounts_manager"
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	dbAccount, _ := sql_db.NewConnection(connAccountString)

	gridRowDataRepositoryWithChild := treegrid.NewGridRowDataRepositoryWithChild(
		db,
		"user_groups",
		"user_group_lines",
		repository.UserGroupFieldNames,
		repository.UserGroupLineFieldNames,
		100,
		&treegrid.GridRowDataRepositoryWithChildCfg{
			MainCol:                  "code",
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
	userGroupService := service.NewUserGroupService(db, gridRowDataRepositoryWithChild)

	grUserGroupDataUploadRepositoryWithChild := treegrid.NewGridRepository(db, "user_groups",
		"user_group_lines",
		repository.UserGroupFieldNames,
		repository.UserGroupLineFieldUploadNames)

	grUserRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"users",
		repository.UserUploadNames,
		1, // arbitrary
	)

	uploadService := service.NewUploadService(db, grUserGroupDataUploadRepositoryWithChild, grUserRepository)

	secretmanager, err := utils.GetSecretManager()
	if err != nil {
		fmt.Printf("error: %v", err)
		log.Panic(err)
	}
	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository, secretmanager)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: userGroupService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) float64 {
			return float64(userGroupService.GetPageCount(tr))
		},
		AccountManagerService: accountService,
	}

	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
