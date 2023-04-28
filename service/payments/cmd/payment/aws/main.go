package main

import (
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/handler"

	lambda_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/handler/aws/lambda"

	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	sql_connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	// ModuleID is hardcoded as provided in the specification.
	moduleID  = 8
	awsRegion = "eu-central-1"

	authorizationHeader = "Authorization"
)

func main() {
	secretsManager, err := secretsmanager.NewAWSSecretsManager(secretsmanager.AWSConfig{
		Region:       awsRegion,
		MaxCacheSize: secretcache.DefaultMaxCacheSize,
		CacheItemTTL: secretcache.DefaultCacheItemTTL,
	})
	if err != nil {
		log.Panic(err)
	}
	appConfig := config.NewAWSSecretsManagerConfig(secretsManager)
	connString := appConfig.GetDBConnection()

	if connString == "" {
		log.Panic("connection string empty!")
	}

	if _, checkConnectionErr := sql_db.NewConnection(connString); checkConnectionErr != nil {
		log.Panicf("could not connect: %v, %s", checkConnectionErr, connString)
	}

	connectionPool := sql_connection.NewPool()
	defer func() {
		if closeErr := connectionPool.Close(); closeErr != nil {
			log.Println(closeErr)
		}
	}()

	paymentLambdaHandler := lambda_handler.NewLambdaHandler(secretsManager, connectionPool)
	lambdaHandler := handler.LambdaTreeGridHandler{
		LambdaPaths: &handler.LambdaTreeGridPaths{
			PathPageCount: "/data",
			PathPageData:  "/page",
			PathUpload:    "/upload",
			PathCell:      "/cell",
		},
		CallbackUploadDataFunc: paymentLambdaHandler.Handle,
	}

	lambda.Start(lambdaHandler)
}
