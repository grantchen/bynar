package config

import "os"

// localConfig is the implementation of AppConfig for local environment
type localConfig struct {
}

const localDBConnKey = "DB_CONN_KEY"
const localAccountConnKey = "DB_ACCOUNT_CONN_KEY"

// GetDBConnection implements AppConfig
func (*localConfig) GetDBConnection() string {
	return os.Getenv(localDBConnKey)
}

// GetAccountManagementConnection implements AppConfig
func (*localConfig) GetAccountManagementConnection() string {
	return os.Getenv(localAccountConnKey)
}

func NewLocalConfig() AppConfig {
	return &localConfig{}
}
