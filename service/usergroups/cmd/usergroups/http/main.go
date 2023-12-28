package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"

	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/internal/service"
)

func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()
	connString := appConfig.GetAccountManagementConnection()
	connAccountString := appConfig.GetAccountManagementConnection()
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	connString = "root:123456@tcp(localhost:3306)/46542255-9d45-49d5-939d-84bc55b1a938"
	db, err := sqldb.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	dbAccount, _ := sqldb.NewConnection(connAccountString)

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
			QueryChildSuggestion:     repository.QueryChildSuggestion,
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
	grUserGroupRepository := treegrid.NewSimpleGridRowRepository(
		db,
		"user_groups",
		repository.UserGroupFieldNames,
		1, // arbitrary
	)

	uploadService := service.NewUploadService(db, grUserGroupRepository, grUserGroupDataUploadRepositoryWithChild, grUserRepository, "en")

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Panic(err)
	}

	accountRepository := pkgrepository.NewAccountManagerRepository(dbAccount)
	accountService := pkgservice.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	h := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: userGroupService.GetPageData,
		CallBackGetCellDataFunc: userGroupService.GetCellSuggestion,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := userGroupService.GetPageCount(tr)
			return float64(count), err
		},
		AccountManagerService: accountService,
	}

	http.HandleFunc("/upload", h.HTTPHandleUpload)
	http.HandleFunc("/data", h.HTTPHandleGetPageCount)
	http.HandleFunc("/page", h.HTTPHandleGetPageData)
	http.HandleFunc("/cell", h.HTTPHandleCell)

	// server
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
