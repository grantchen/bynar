package main

import (
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/config"
	lambda_handler "git-codecommit.eu-central-1.amazonaws.com/v1/repos/payments/internal/handler/aws/lambda"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"

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

	lambdaHandler := lambda_handler.NewLambdaHandler(secretsManager, connectionPool)
	lambda.Start(lambdaHandler)
}
