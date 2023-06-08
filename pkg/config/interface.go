package config

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"

type AppConfig interface {
	GetDBConnection() string
	GetAccountManagementConnection() string
}

func GetAppConfig(s secretsmanager.SecretsManager) AppConfig {
	return NewAWSSecretsManagerConfig(s)
}
