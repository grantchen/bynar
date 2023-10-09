package main

import (
	"fmt"
	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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
	general_posting_setup_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/external/service"
	invoices_service "git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/external/handler/service"
)

type HandlerMapping struct {
	handler    *handler.HTTPTreeGridHandler
	prefixPath string
}

type HandlerMappingWithPermission struct {
	factoryFunc treegrid.TreeGridServiceFactoryFunc
	prefixPath  string
}

const prefix = "/apprunnerurl"

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file in main service ", err)
	}
	appConfig := config.NewLocalConfig()

	connectionPool := connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	connString := appConfig.GetDBConnection()
	db, err := sql_db.InitializeConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	logger.Debug("connection string account: ", accountManagementConnectionString)

	accountDB, err := sql_db.NewConnection(accountManagementConnectionString)
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
	http.Handle("/user", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.User))))
	http.Handle("/user/", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.GetUserProfileById))))
	http.Handle("/upload", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.UploadProfilePhoto))))
	http.Handle("/profile-image", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.DeleteProfileImage))))
	http.Handle("/update-user-language-preference", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.UpdateUserLanguagePreference))))
	http.Handle("/update-user-theme-preference", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.UpdateUserThemePreference))))
	http.Handle("/update-profile", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.UpdateUserProfile))))

	// TreeGrid Components that require authentication, but not permission
	accountRepository := pkg_repository.NewAccountManagerRepository(accountDB)
	accountService := pkg_service.NewAccountManagerService(accountDB, accountRepository, authProvider)
	lsHandlerMappingWithAuthentication := make([]*HandlerMappingWithPermission, 0)
	lsHandlerMappingWithAuthentication = append(lsHandlerMappingWithAuthentication,
		&HandlerMappingWithPermission{factoryFunc: invoices_service.NewTreeGridServiceFactory(), prefixPath: "/invoices"},
	)
	for _, handlerMappingWithAuthentication := range lsHandlerMappingWithAuthentication {
		handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
			AccountManagerService:  accountService,
			TreeGridServiceFactory: handlerMappingWithAuthentication.factoryFunc,
			ConnectionPool:         connectionPool,
			PathPrefix:             prefix + handlerMappingWithAuthentication.prefixPath,
			IsValidatePermissions:  false,
		}
		handler.HandleHTTPReqWithAuthenMWAndDefaultPath()
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

	for _, handlerMapping := range lsHandlerMapping {
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/data", handlerMapping.handler.HTTPHandleGetPageCount)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/page", handlerMapping.handler.HTTPHandleGetPageData)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/upload", handlerMapping.handler.HTTPHandleUpload)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/cell", handlerMapping.handler.HTTPHandleCell)

	}

	// _______________________________________ COMPONENT WITH PERMISSION______________________________________________
	lsHandlerMappingWithPermission := make([]*HandlerMappingWithPermission, 0)
	lsHandlerMappingWithPermission = append(lsHandlerMappingWithPermission,
		&HandlerMappingWithPermission{factoryFunc: organizations_service.NewTreeGridServiceFactory(), prefixPath: "/organizations"},
		&HandlerMappingWithPermission{factoryFunc: accounts_service.NewTreeGridServiceFactory(), prefixPath: "/user_list"},
		&HandlerMappingWithPermission{factoryFunc: general_posting_setup_service.NewTreeGridServiceFactory(), prefixPath: "/general_posting_setup"},
	)

	for _, handlerMappingWithPermission := range lsHandlerMappingWithPermission {
		handler := &handler.HTTPTreeGridHandlerWithDynamicDB{
			AccountManagerService:  accountService,
			TreeGridServiceFactory: handlerMappingWithPermission.factoryFunc,
			ConnectionPool:         connectionPool,
			PathPrefix:             prefix + handlerMappingWithPermission.prefixPath,
			IsValidatePermissions:  true,
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
