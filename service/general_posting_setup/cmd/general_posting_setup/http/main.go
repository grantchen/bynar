package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/internal/service"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

func main() {
	connString := "root:123456@tcp(localhost:3306)/bynar"
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sqldb.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleGeneralPostingSetupRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(
		db,
		"general_posting_setup",
		repository.GeneralPostingSetupFieldNames,
		100,
		&treegrid.SimpleGridRepositoryCfg{
			MainCol:     "code",
			QueryString: repository.QuerySelect,
			QueryJoin:   repository.QueryJoin,
			QueryCount:  repository.QueryCount,
		},
	)

	generalPostingSetupRepository := repository.NewPostingSetupRepository(db)
	generalPostingSetupService := service.NewGeneralPostingSetupService(simpleGeneralPostingSetupRepository)
	uploadService := service.NewUploadService(
		db,
		simpleGeneralPostingSetupRepository,
		generalPostingSetupRepository,
		"",
	)

	h := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: generalPostingSetupService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := generalPostingSetupService.GetPageCount(tr)
			return float64(count), err
		},
	}
	http.HandleFunc("/apprunnerurl/general_posting_setup/upload", h.HTTPHandleUpload)
	http.HandleFunc("/apprunnerurl/general_posting_setup/data", h.HTTPHandleGetPageCount)
	http.HandleFunc("/apprunnerurl/general_posting_setup/page", h.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
