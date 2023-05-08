package utils

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	"github.com/aws/aws-secretsmanager-caching-go/secretcache"
)

var (
	// ModuleID is hardcoded as provided in the specification.
	awsRegion = "eu-central-1"
)

func GetSecretManager() (secretsmanager.SecretsManager, error) {
	secretsManager, err := secretsmanager.NewAWSSecretsManager(secretsmanager.AWSConfig{
		Region:       awsRegion,
		MaxCacheSize: secretcache.DefaultMaxCacheSize,
		CacheItemTTL: secretcache.DefaultCacheItemTTL,
	})

	return secretsManager, err
}
