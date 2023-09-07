package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/http_handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"github.com/joho/godotenv"
)

// use for test only this module without permission
func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()

	db, err := sql_db.NewConnection(appConfig.GetAccountManagementConnection())
	if err != nil {
		log.Fatal(err)
	}

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Fatal(err)
	}
	paymentProvider, err := checkout.NewPaymentClient()
	if err != nil {
		log.Fatal(err)
	}

	handler := http_handler.NewHTTPHandler(db, authProvider, paymentProvider)

	// Signup endpoints
	http.HandleFunc("/signup", handler.Signup)
	http.HandleFunc("/confirm-email", handler.ConfirmEmail)
	http.HandleFunc("/verify-card", handler.VerifyCard)
	http.HandleFunc("/create-user", handler.CreateUser)

	// Signin endpoints
	http.HandleFunc("/signin", handler.Signin)

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
