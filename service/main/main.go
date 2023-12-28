package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gcs"

	"github.com/joho/godotenv"

	accountshttphandler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/http"
	accountsservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/accounts/external/handler/service"
	cardshttphandler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/cards/external/handler/http"
	generalpostingsetupservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/general_posting_setup/external/service"
	invoicesservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/invoices/external/handler/service"
	languagesservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/languages/external/handler/service"
	organizationsservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/external/handler/service"
	paymentsservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/external/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/checkout"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	sqldb "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/gip"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/render"
	pkgrepository "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/repository"
	pkgservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/service"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/treegrid"
	procurementsservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/external/service"
	saleservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/external/service"
	sitesservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/sites/external/handler/service"
	transfersservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/external/handler/service"
	usergroupservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/usergroups/external/service"
	warehousesservice "git-codecommit.eu-central-1.amazonaws.com/v1/repos/warehouses/external/service"
)

type HandlerMapping struct {
	handler    *handler.HTTPTreeGridHandler
	prefixPath string
}

type HandlerMappingWithPermission struct {
	factoryFunc treegrid.ServiceFactoryFunc
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

	accountManagementConnectionString := appConfig.GetAccountManagementConnection()
	logger.Debug("connection string account: ", accountManagementConnectionString)

	accountDB, err := sqldb.NewConnection(accountManagementConnectionString)
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

	accountHandler := accountshttphandler.NewHTTPHandler(accountDB, authProvider, paymentProvider, cloudStorageProvider)

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
	http.Handle("/organization-account", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.GetOrganizationAccount))))
	http.Handle("/update-organization-account", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.UpdateOrganizationAccount))))
	http.Handle("/delete-organization-account", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(accountHandler.DeleteOrganizationAccount))))

	cardsHandler := cardshttphandler.NewHTTPHandler(accountDB, authProvider, paymentProvider)
	// Cards endpoints
	http.Handle(prefix+"/cards/list", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(cardsHandler.ListCards))))
	http.Handle(prefix+"/cards/add", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(cardsHandler.AddCard))))
	http.Handle(prefix+"/cards/update", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(cardsHandler.UpdateCard))))
	http.Handle(prefix+"/cards/delete", render.CorsMiddleware(handler.VerifyIdTokenAndInitDynamicDB(http.HandlerFunc(cardsHandler.DeleteCard))))

	// TreeGrid Components that require authentication, but not permission
	accountRepository := pkgrepository.NewAccountManagerRepository(accountDB)
	accountService := pkgservice.NewAccountManagerService(accountDB, accountRepository, authProvider)
	lsHandlerMappingWithAuthentication := make([]*HandlerMappingWithPermission, 0)
	lsHandlerMappingWithAuthentication = append(lsHandlerMappingWithAuthentication,
		&HandlerMappingWithPermission{factoryFunc: invoicesservice.NewTreeGridServiceFactory(), prefixPath: "/invoices"},
	)
	for _, handlerMappingWithAuthentication := range lsHandlerMappingWithAuthentication {
		h := &handler.HTTPTreeGridHandlerWithDynamicDB{
			AccountManagerService:  accountService,
			TreeGridServiceFactory: handlerMappingWithAuthentication.factoryFunc,
			ConnectionPool:         connectionPool,
			PathPrefix:             prefix + handlerMappingWithAuthentication.prefixPath,
			IsValidatePermissions:  false,
		}
		h.HandleHTTPReqWithAuthenMWAndDefaultPath()
	}

	lsHandlerMapping := make([]*HandlerMapping, 0)
	//lsHandlerMapping = append(lsHandlerMapping,
	//	&HandlerMapping{handler: payments_handler.NewHTTPHandler(appConfig, db),
	//		prefixPath: "/payments"})

	for _, handlerMapping := range lsHandlerMapping {
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/data", handlerMapping.handler.HTTPHandleGetPageCount)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/page", handlerMapping.handler.HTTPHandleGetPageData)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/upload", handlerMapping.handler.HTTPHandleUpload)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/cell", handlerMapping.handler.HTTPHandleCell)

	}

	// _______________________________________ COMPONENT WITH PERMISSION______________________________________________
	lsHandlerMappingWithPermission := make([]*HandlerMappingWithPermission, 0)
	lsHandlerMappingWithPermission = append(lsHandlerMappingWithPermission,
		&HandlerMappingWithPermission{factoryFunc: organizationsservice.NewTreeGridServiceFactory(), prefixPath: "/organizations"},
		&HandlerMappingWithPermission{factoryFunc: sitesservice.NewTreeGridServiceFactory(), prefixPath: "/sites"},
		&HandlerMappingWithPermission{factoryFunc: transfersservice.NewTreeGridServiceFactory(), prefixPath: "/transfers"},
		&HandlerMappingWithPermission{factoryFunc: accountsservice.NewTreeGridServiceFactory(), prefixPath: "/user_list"},
		&HandlerMappingWithPermission{factoryFunc: generalpostingsetupservice.NewTreeGridServiceFactory(), prefixPath: "/general_posting_setup"},
		&HandlerMappingWithPermission{factoryFunc: usergroupservice.NewTreeGridServiceFactory(), prefixPath: "/user_groups"},
		&HandlerMappingWithPermission{factoryFunc: warehousesservice.NewTreeGridServiceFactory(), prefixPath: "/warehouses"},
		&HandlerMappingWithPermission{factoryFunc: saleservice.NewTreeGridServiceFactory(), prefixPath: "/sales"},
		&HandlerMappingWithPermission{factoryFunc: paymentsservice.NewTreeGridServiceFactory(), prefixPath: "/payments"},
		&HandlerMappingWithPermission{factoryFunc: procurementsservice.NewTreeGridServiceFactory(), prefixPath: "/procurements"},
		&HandlerMappingWithPermission{factoryFunc: languagesservice.NewTreeGridServiceFactory(), prefixPath: "/languages"},
	)

	for _, handlerMappingWithPermission := range lsHandlerMappingWithPermission {
		h := &handler.HTTPTreeGridHandlerWithDynamicDB{
			AccountManagerService:  accountService,
			TreeGridServiceFactory: handlerMappingWithPermission.factoryFunc,
			ConnectionPool:         connectionPool,
			PathPrefix:             prefix + handlerMappingWithPermission.prefixPath,
			IsValidatePermissions:  true,
		}
		h.HandleHTTPReqWithAuthenMWAndDefaultPath()
	}

	log.Println("start server at 8080!")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
