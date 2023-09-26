package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/internal/service"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// use for test only this module without permission
func main() {
	// gr := &treegrid.GridRow{}
	// (*gr)["id"] = "1"
	// (*gr)["name"] = "linh"
	// (*gr)["vat_number"] = "123"
	// (*gr)["state"] = "1123"
	// (*gr)["code"] = "abc"
	// s, _ := gr.MakeUpdateQuery("invoices", repository.InvoiceFieldNames)
	// fmt.Println(s)

	// secretmanager, err := utils.GetSecretManager()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// appConfig := config.NewAWSSecretsManagerConfig(secretmanager)
	// connString := appConfig.GetDBConnection()
	connString := "root:123456@tcp(localhost:3306)/bynar"
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleInvoiceRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "invoices", repository.InvoiceFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:     "code",
			QueryString: repository.QuerySelect,
			QueryCount:  repository.QueryCount,
		})

	uploadService, _ := service.NewUploadService(db, simpleInvoiceRepository, 0)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: uploadService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := uploadService.GetPageCount(tr)
			return float64(count), err
		},
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
