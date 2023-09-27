package main

import (
	"fmt"
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
	connString := "root:123456@tcp(localhost:3306)/bynar"
	db, err := sql_db.NewConnection(connString)
	if err != nil {
		log.Panic(err)
	}

	// TODO
	// test data
	accountID := 1
	simpleInvoiceRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "invoices", repository.InvoiceFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:       "code",
			QueryString:   repository.QuerySelect,
			QueryCount:    repository.QueryCount,
			AdditionWhere: fmt.Sprintf(repository.AdditionWhere, accountID),
		})
	treeGridService, _ := service.NewTreeGridService(db, simpleInvoiceRepository, accountID)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:   treeGridService.Handle,
		CallbackGetPageDataFunc:  treeGridService.GetPageData,
		CallbackGetPageCountFunc: treeGridService.GetPageCount,
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
