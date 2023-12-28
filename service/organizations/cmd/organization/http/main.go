package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// use for test only this module without permission
func main() {
	// connString := appConfig.GetDBConnection()
	connString := "root:123456@tcp(localhost:3306)/bynar"
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sqldb.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "organizations", repository.OrganizationFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:     "code",
			QueryString: repository.QuerySelect,
			QueryCount:  repository.QueryCount,
		})
	organizationService := service.NewOrganizationService(db, simpleOrganizationRepository)

	userID := 1 // for test only
	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository, userID, "en")

	h := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: organizationService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := organizationService.GetPageCount(tr)
			return float64(count), err
		},
	}
	http.HandleFunc("/upload", h.HTTPHandleUpload)
	http.HandleFunc("/data", h.HTTPHandleGetPageCount)
	http.HandleFunc("/page", h.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
