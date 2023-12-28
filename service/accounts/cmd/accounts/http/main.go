package main

import (
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
)

const prefix = "/apprunnerurl"

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
	gcsProvider, err := gcs.NewGCSClient()
	if err != nil {
		log.Fatal(err)
	}

	accountRepository := pkgrepository.NewAccountManagerRepository(dbAccount)
	accountService := pkgservice.NewAccountManagerService(dbAccount, accountRepository, authProvider)

	httpHandler := http_handler.NewHTTPHandler(dbAccount, authProvider, paymentProvider, gcsProvider)

	// Signup endpoints
	http.Handle("/signup", render.CorsMiddleware(http.HandlerFunc(httpHandler.Signup)))
	http.Handle("/confirm-email", render.CorsMiddleware(http.HandlerFunc(httpHandler.ConfirmEmail)))
	http.Handle("/verify-card", render.CorsMiddleware(http.HandlerFunc(httpHandler.VerifyCard)))
	http.Handle("/create-user", render.CorsMiddleware(http.HandlerFunc(httpHandler.CreateUser)))

	// Signin endpoints
	http.Handle("/signin-email", render.CorsMiddleware(http.HandlerFunc(httpHandler.SendSignInEmail)))
	http.Handle("/signin", render.CorsMiddleware(http.HandlerFunc(httpHandler.SignIn)))

	// user endpoints
	http.Handle("/user", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.User))))
	http.Handle("/upload", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.UploadProfilePhoto))))
	http.Handle("/profile-image", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.DeleteProfileImage))))
	http.Handle("/update-user-language-preference", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.UpdateUserLanguagePreference))))
	http.Handle("/update-user-theme-preference", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.UpdateUserThemePreference))))
	http.Handle("/update-profile", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(httpHandler.UpdateUserProfile))))

	// accounts treegrid endpoints
	dbHandler := &handler.HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountService,
		TreeGridServiceFactory: service.NewTreeGridServiceFactory(),
		ConnectionPool:         connectionPool,
		PathPrefix:             prefix + "/user_list",
	}
	dbHandler.HandleHTTPReqWithAuthenMWAndDefaultPath()

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
