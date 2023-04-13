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

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
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

	db, err := sql_db.NewConnection(connString)
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
	// port := os.Getenv("port")

	// if port == "" {
	// 	port = "9002"
	// }
	// log.Fatal(http.ListenAndServe(":9002", myRouter))

	log.Println("Listening on :9003...")
	err := http.ListenAndServe(":9003", nil)
	if err != nil {
		log.Fatal(err)
	}
}
