package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/service"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	handleRequests()

}

func initMySQLDB() (*sql.DB, error) {
	connString := config.GetMySqlUsername() + ":" + config.GetMySqlPassword() + "@tcp(" +
		config.GetMySqlHost() + ":" + config.GetMySqlPort() + ")/" + config.GetMySqlDBName()

	fmt.Println(connString)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}
	// defer db.Close()
	return db, err
}

func handleRequests() {
	config.Load()

	mysqlDB, _ := initMySQLDB()

	dataGridRepo := repository.NewMysqlDataGridRepository(mysqlDB)
	dataGridService := service.NewDataGridService(dataGridRepo)
	httpDataGridHandler := handler.NewHTTPDataGridHandler(dataGridService)

	http.HandleFunc("/TreeGridView", httpDataGridHandler.GetTreeGridViewData)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fs2 := http.FileServer(http.Dir("./Grid"))
	http.Handle("/Grid/", http.StripPrefix("/Grid/", fs2))

	log.Println("Listening on :9003...")
	err := http.ListenAndServe(":9003", nil)
	if err != nil {
		log.Fatal(err)
	}
}
