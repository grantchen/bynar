package main

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"log"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const prefix = "/apprunnerurl"

type HandlerMappingWithDynamicDB struct {
	Path        string
	RequestFunc func(w http.ResponseWriter, r *http.Request)
}

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
	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository)

	// db, err := sql_db.NewConnection(appConfig.GetAccountManagementConnection())
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
	lsHandlerMappingWithDynamicDB := make([]*HandlerMappingWithDynamicDB, 0)
	lsHandlerMappingWithDynamicDB = append(lsHandlerMappingWithDynamicDB,
		&HandlerMappingWithDynamicDB{Path: "/user", RequestFunc: httpHandler.User},
		&HandlerMappingWithDynamicDB{Path: "/upload", RequestFunc: httpHandler.UploadProfilePhoto},
		&HandlerMappingWithDynamicDB{Path: "/profile-image", RequestFunc: httpHandler.DeleteProfileImage},
	)
	for _, handlerMappingWithPermission := range lsHandlerMappingWithDynamicDB {
		handler := &handler.HTTPHandlerWithDynamicDB{
			ConnectionPool: connectionPool,
			Path:           handlerMappingWithPermission.Path,
			RequestFunc:    handlerMappingWithPermission.RequestFunc,
		}
		handler.HandleHTTPReqWithDynamicDB()
	}

	// accounts treegrid endpoints
	dbhandler := &handler.HTTPTreeGridHandlerWithDynamicDB{
		AccountManagerService:  accountService,
		TreeGridServiceFactory: service.NewTreeGridServiceFactory(),
		ConnectionPool:         connectionPool,
		PathPrefix:             prefix + "/user_list",
	}
	dbhandler.HandleHTTPReqWithAuthenMWAndDefaultPath()

	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
