package main

import (
	"log"
	"net/http"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/internal/service"
)

func main() {
	connString := "root:123456@tcp(localhost:3306)/bynar"
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleWarehousesRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(
		db,
		"warehouses",
		repository.WarehousesFieldNames,
		100,
		&treegrid.SimpleGridRepositoryCfg{
			MainCol:     "code",
			QueryString: repository.QuerySelect,
			QueryJoin:   repository.QueryJoin,
			QueryCount:  repository.QueryCount,
		},
	)

	warehousesRepository := repository.NewPostingSetupRepository(db)
	warehousesService := service.NewWarehousesService(simpleWarehousesRepository)
	uploadService := service.NewUploadService(
		db,
		simpleWarehousesRepository,
		warehousesRepository,
		"en",
	)

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: warehousesService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := warehousesService.GetPageCount(tr)
			return float64(count), err
		},
	}
	http.HandleFunc("/apprunnerurl/warehouses/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/apprunnerurl/warehouses/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/apprunnerurl/warehouses/page", handler.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
