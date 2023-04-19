package main

import (
	"log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	sql_db "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/transfers/internal/config"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	// ModuleID is hardcoded as provided in the specification.
	moduleID = 8

	connectionStringKey = "bynar"
	awsRegion           = "eu-central-1"

	authorizationHeader = "Authorization"
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

}
