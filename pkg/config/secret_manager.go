package config

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/aws/secretsmanager"
	sql_connection "git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/db/connection"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/pkgs/logger"
)

type awsSecretAppConfig struct {
	secretsmanager secretsmanager.SecretsManager
}

const dbConnKey = "bynar"

// GetDBConnection implements AppConfig
func (a *awsSecretAppConfig) GetDBConnection() string {
	value, err := a.secretsmanager.GetString(dbConnKey)

	if err != nil {
		logger.Debug(err)
		return ""

	}
	return sql_connection.JSON2DatabaseConnection(value)
}

func NewAWSSecretsManagerConfig(secretsManager secretsmanager.SecretsManager) AppConfig {
	return &awsSecretAppConfig{secretsmanager: secretsManager}
}
