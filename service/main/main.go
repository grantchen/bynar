package main

import (
	"fmt"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	general_posting_setup_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/external/handler/http"
	organizations_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/external/handler/service"
	payments_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/external/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	pkg_repository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkg_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	procurements_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/external/handler/http"
	sales_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/external/handler/http"
	transfers_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/external/handler/http"
	usergroups_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/external/handler/http"

	accounts_http_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/http"
	accounts_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/service"
)

type HandlerMapping struct {
	handler    *handler.HTTPTreeGridHandler
	prefixPath string
}

type HandlerMappingWithPermission struct {
	factoryFunc treegrid.TreeGridServiceFactoryFunc
	prefixPath  string
}

type HandlerMappingWithDynamicDB struct {
	Path        string
	RequestFunc func(w http.ResponseWriter, r *http.Request)
}

const prefix = "/apprunnerurl"

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file in main service")
	}

	appConfig := config.NewLocalConfig()

	connectionPool := connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	connString := appConfig.GetDBConnection()
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	accountDB, err := sql_db.NewConnection(appConfig.GetAccountManagementConnection())
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

	cloudStorageProvider, err := gcs.NewGCSClient()
	if err != nil {
		log.Fatal(err)
	}

	accountHandler := accounts_http_handler.NewHTTPHandler(accountDB, authProvider, paymentProvider, cloudStorageProvider)

	// Signup endpoints
	http.Handle("/signup", render.CorsMiddleware(http.HandlerFunc(accountHandler.Signup)))
	http.Handle("/confirm-email", render.CorsMiddleware(http.HandlerFunc(accountHandler.ConfirmEmail)))
	http.Handle("/verify-card", render.CorsMiddleware(http.HandlerFunc(accountHandler.VerifyCard)))
	http.Handle("/create-user", render.CorsMiddleware(http.HandlerFunc(accountHandler.CreateUser)))

	// Signin endpoints
	http.Handle("/signin-email", render.CorsMiddleware(http.HandlerFunc(accountHandler.SendSignInEmail)))
	http.Handle("/signin", render.CorsMiddleware(http.HandlerFunc(accountHandler.SignIn)))

	// user endpoints
	lsHandlerMappingWithDynamicDB := make([]*HandlerMappingWithDynamicDB, 0)
	lsHandlerMappingWithDynamicDB = append(lsHandlerMappingWithDynamicDB,
		&HandlerMappingWithDynamicDB{Path: "/user", RequestFunc: accountHandler.User},
		&HandlerMappingWithDynamicDB{Path: "/upload", RequestFunc: accountHandler.UploadProfilePhoto},
		&HandlerMappingWithDynamicDB{Path: "/profile-image", RequestFunc: accountHandler.DeleteProfileImage},
		&HandlerMappingWithDynamicDB{Path: "/update-user-language-preference", RequestFunc: accountHandler.UpdateUserLanguagePreference},
		&HandlerMappingWithDynamicDB{Path: "/update-user-theme-preference", RequestFunc: accountHandler.UpdateUserThemePreference},
	)
	for _, handlerMappingWithPermission := range lsHandlerMappingWithDynamicDB {
		handler := &handler.HTTPHandlerWithDynamicDB{
			ConnectionPool: connectionPool,
			Path:           handlerMappingWithPermission.Path,
			RequestFunc:    handlerMappingWithPermission.RequestFunc,
		}
		handler.HandleHTTPReqWithDynamicDB()
	}

	lsHandlerMapping := make([]*HandlerMapping, 0)
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: sales_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/sales"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: payments_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/payments"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: transfers_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/transfers"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: procurements_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/procurements"})
	// lsHandlerMapping = append(lsHandlerMapping,
	// 	&HandlerMapping{handler: organizations_handler.NewHTTPHandler(appConfig, db),
	// 		prefixPath: "/organizations"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: usergroups_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/user_groups"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: general_posting_setup_handler.NewHTTPHandler(appConfig, db),
			prefixPath: "/general_posting_setup"})

	for _, handlerMapping := range lsHandlerMapping {
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/data", handlerMapping.handler.HTTPHandleGetPageCount)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/page", handlerMapping.handler.HTTPHandleGetPageData)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/upload", handlerMapping.handler.HTTPHandleUpload)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/cell", handlerMapping.handler.HTTPHandleCell)

	}

	// _______________________________________ COMPONENT WITH PERMISSION______________________________________________
	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	logger.Debug("connection string account: ", accountManagementConnectionString)

	dbAccount, err := connectionPool.Get(accountManagementConnectionString)

	if err != nil {
		log.Panic(err)
	}
	accountRepository := pkg_repository.NewAccountManagerRepository(dbAccount)
	accountService := pkg_service.NewAccountManagerService(dbAccount, accountRepository)

	lsHandlerMappingWithPermission := make([]*HandlerMappingWithPermission, 0)
	lsHandlerMappingWithPermission = append(lsHandlerMappingWithPermission,
		&HandlerMappingWithPermission{factoryFunc: organizations_service.NewTreeGridServiceFactory(), prefixPath: "/organizations"},
		&HandlerMappingWithPermission{factoryFunc: accounts_service.NewTreeGridServiceFactory(), prefixPath: "/user_list"},
	)

	for _, handlerMappingWithPermission := range lsHandlerMappingWithPermission {
		handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
			AccountManagerService:  accountService,
			TreeGridServiceFactory: handlerMappingWithPermission.factoryFunc,
			ConnectionPool:         connectionPool,
			PathPrefix:             prefix + handlerMappingWithPermission.prefixPath,
		}
		handler.HandleHTTPReqWithAuthenMWAndDefaultPath()
	}

	log.Println("start server at 8080!")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
