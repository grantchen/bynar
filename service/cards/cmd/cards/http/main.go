package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/external/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const prefix = "/apprunnerurl/cards"

// use for test only this module without permission
func main() {
	err := godotenv.Load("../main/.env")
	if err != nil {
		log.Fatal("Error loading .env file in main service")
	}
	appConfig := config.NewLocalConfig()

	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	logrus.Debug("connection string account: ", accountManagementConnectionString)

	connectionPool := connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()
	dbAccount, err := connectionPool.Get(accountManagementConnectionString)

	if err != nil {
		log.Panic(err)
	}

	authProvider, err := gip.NewGIPClient()
	if err != nil {
		log.Fatal(err)
	}
	paymentProvider, err := checkout.NewPaymentClient()
	if err != nil {
		log.Fatal(err)
	}

	httpHandler := http_handler.NewHTTPHandler(dbAccount, authProvider, paymentProvider)

	// Cards endpoints
	http.Handle(prefix+"/list", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.ListCards))))
	http.Handle(prefix+"/add", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.AddCard))))
	http.Handle(prefix+"/update", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.UpdateCard))))
	http.Handle(prefix+"/delete", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.DeleteCard))))

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
