package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/http_handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
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
	http.Handle("/signup", render.CorsMiddleware(http.HandlerFunc(handler.Signup)))
	http.Handle("/confirm-email", render.CorsMiddleware(http.HandlerFunc(handler.ConfirmEmail)))
	http.Handle("/verify-card", render.CorsMiddleware(http.HandlerFunc(handler.VerifyCard)))
	http.Handle("/create-user", render.CorsMiddleware(http.HandlerFunc(handler.CreateUser)))

	// Signin endpoints
	http.Handle("/signin-email", render.CorsMiddleware(http.HandlerFunc(handler.SendSignInEmail)))
	http.Handle("/signin", render.CorsMiddleware(http.HandlerFunc(handler.SignIn)))
	// user endpoints
	http.Handle("/user", render.CorsMiddleware(http.HandlerFunc(handler.User)))

	// TreeGrid handler endpoints
	tgHandler := http_handler.NewUserHTTPHandler()

	http.Handle("/upload", render.CorsMiddleware(http.HandlerFunc(tgHandler.HTTPHandleUpload)))
	http.Handle("/data", render.CorsMiddleware(http.HandlerFunc(tgHandler.HTTPHandleGetPageCount)))
	http.Handle("/page", render.CorsMiddleware(http.HandlerFunc(tgHandler.HTTPHandleGetPageData)))

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
