package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/internal/service"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
)

// use for test only this module without permission
func main() {
	// connString := appConfig.GetDBConnection()
	connString := "root:123456@tcp(localhost:3306)/bynar"
	// connString := "root:Munrfe2020@tcp(bynar-cet.ccwuyxj7ucnd.eu-central-1.rds.amazonaws.com:3306)/bynar"
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	simpleLanguageRepository := treegrid.NewSimpleGridRowRepositoryWithCfg(db, "languages", repository.LanguageFieldNames,
		100, &treegrid.SimpleGridRepositoryCfg{
			MainCol:     "language",
			QueryString: repository.QuerySelect,
			QueryCount:  repository.QueryCount,
		})
	languageService := service.NewLanguageService(db, simpleLanguageRepository)

	uploadService, _ := service.NewUploadService(db, languageService, simpleLanguageRepository, "en")

	handler := &handler.HTTPTreeGridHandler{
		CallbackUploadDataFunc:  uploadService.Handle,
		CallbackGetPageDataFunc: languageService.GetPageData,
		CallbackGetPageCountFunc: func(tr *treegrid.Treegrid) (float64, error) {
			count, err := languageService.GetPageCount(tr)
			return float64(count), err
		},
	}
	http.HandleFunc("/upload", handler.HTTPHandleUpload)
	http.HandleFunc("/data", handler.HTTPHandleGetPageCount)
	http.HandleFunc("/page", handler.HTTPHandleGetPageData)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
