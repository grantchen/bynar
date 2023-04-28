package main

import (
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/repository"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	// ModuleID is hardcoded as provided in the specification.
	awsRegion = "eu-central-1"
)

func getAppConfig(s secretsmanager.SecretsManager) config.AppConfig {
	return config.NewAWSSecretsManagerConfig(s)
}

func main() {
	secretsManager, err := secretsmanager.NewAWSSecretsManager(secretsmanager.AWSConfig{
		Region:       awsRegion,
		MaxCacheSize: secretcache.DefaultMaxCacheSize,
		CacheItemTTL: secretcache.DefaultCacheItemTTL,
	})

	if err != nil {
		log.Panic(err)
	}
	appConfig := getAppConfig(secretsManager)
	connString := appConfig.GetDBConnection()
	db, err := sql_db.NewConnection(connString)

	if err != nil {
		log.Panic(err)
	}

	transferRepository := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(transferRepository)
	// handler := lambda_handler.NewLambdaHandler(transferService)

	lambdaHandler := handler.LambdaTreeGridHandler{
		LambdaPaths: &handler.LambdaTreeGridPaths{
			PathPageCount: "/data",
			PathPageData:  "/page",
			PathUpload:    "/upload",
			PathCell:      "/cell",
		},
		CallbackGetPageCountFunc: transferService.GetPagesCount,
		CallbackGetPageDataFunc:  transferService.GetTransfersPageData,
	}

	lambda.Start(lambdaHandler)
}
