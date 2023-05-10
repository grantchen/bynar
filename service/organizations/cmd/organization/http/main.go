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
	connString := "root:123456@tcp(localhost:3306)/bynar"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleOrganizationRepository := treegrid.NewSimpleGridRowRepository(db, "organizations", repository.OrganizationFieldNames)
	organizationService := service.NewOrganizationService(db)

	uploadService, _ := service.NewUploadService(db, organizationService, simpleOrganizationRepository)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc: uploadService.Handle,
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
