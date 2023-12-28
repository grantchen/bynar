package main

import (
	"database/sql"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cellurls/internal/service"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"

	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
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

	db, err := sqldb.NewConnection(connString)
	if err != nil {
		panic(err.Error())
	}
	// defer db.Close()
	return db, err
}

func handleRequests() {
	_ = config.Load()

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
