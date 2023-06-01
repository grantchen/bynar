package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/internal/service"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func main() {
	// gr := &treegrid.GridRow{}
	// (*gr)["id"] = "1"
	// (*gr)["name"] = "linh"
	// (*gr)["vat_number"] = "123"
	// (*gr)["state"] = "1123"
	// (*gr)["code"] = "abc"
	// s, _ := gr.MakeUpdateQuery("organizations", repository.OrganizationFieldNames)
	// fmt.Println(s)

	// secretmanager, err := utils.GetSecretManager()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// appConfig := config.NewAWSSecretsManagerConfig(secretmanager)
	// connString := appConfig.GetDBConnection()
	// connString := "root:123456@tcp(localhost:3306)/bynar"
	connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "organizations", repository.OrganizationFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{MainCol: "code"})
	organizationService := service.NewOrganizationService(db, simpleOrganizationRepository)

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: organizationService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) float64 {
			return float64(organizationService.GetPageCount(tr))
		},
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
