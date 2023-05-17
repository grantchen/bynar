package main

import (
	"fmt"
	"log"
	"net/http"

	organizations_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/organizations/external/handler/http"
	payments_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/external/handler/http"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/utils"
	procurements_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/procurements/external/handler/http"
	sales_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/sales/external/handler/http"
	transfers_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/external/handler/http"
)

type HandlerMapping struct {
	handler    *handler.HTTPTreeGridHandler
	prefixPath string
}

const prefix = "/apprunnerurl"

func main() {

	secretmanager, err := utils.GetSecretManager()
	if err != nil {
		fmt.Printf("error: %v", err)
		log.Panic(err)
	}

	appConfig := config.NewAWSSecretsManagerConfig(secretmanager)

	lsHandlerMapping := make([]*HandlerMapping, 0)
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: sales_handler.NewHTTPHandler(appConfig),
			prefixPath: "/sales"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: payments_handler.NewHTTPHandler(appConfig),
			prefixPath: "/payments"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: transfers_handler.NewHTTPHandler(appConfig),
			prefixPath: "/transfers"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: procurements_handler.NewHTTPHandler(appConfig),
			prefixPath: "/procurements"})
	lsHandlerMapping = append(lsHandlerMapping,
		&HandlerMapping{handler: organizations_handler.NewHTTPHandler(appConfig),
			prefixPath: "/organizations"})

	for _, handlerMapping := range lsHandlerMapping {
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/data", handlerMapping.handler.HTTPHandleGetPageCount)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/page", handlerMapping.handler.HTTPHandleGetPageData)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/upload", handlerMapping.handler.HTTPHandleUpload)
		http.HandleFunc(prefix+handlerMapping.prefixPath+"/cell", handlerMapping.handler.HandleCell)

	}
	log.Println("start server at 8080!")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
