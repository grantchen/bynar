package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	httphandler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/simplecurd/internal/service"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	handleRequests()

}

func initMySQLDB() (*sql.DB, error) {
	var err error
	cfg := mysql.Config{}
	file, err := os.Open("./config/config.development.json")
	if err != nil {
		log.Fatalln(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sqldb.NewConnection(cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}
	// defer db.Close()
	return db, err
}

func handleRequests() {

	mysqlDB, _ := initMySQLDB()

	languageRepostory := repository.NewLanguageRepository(mysqlDB, "languages")
	languageService := service.NewLanguageService(languageRepostory)
	httpLanguageHandler := httphandler.NewLanguageHandler(languageService)

	http.HandleFunc("/languages", httpLanguageHandler.HandleGetAllLanguage)
	http.HandleFunc("/languages/update", httpLanguageHandler.HandleUpdateRequest)

	log.Println("Listening on :9003...")
	err := http.ListenAndServe(":9003", nil)
	if err != nil {
		log.Fatal(err)
	}
}
