package config

import "os"

type localConfig struct {
}

const localDBConnKey = "DB_CONN_KEY"

// GetDBConnection implements AppConfig
func (*localConfig) GetDBConnection() string {
	return os.Getenv(localDBConnKey)
}

func NewLocalConfig() AppConfig {
	return &localConfig{}
}
