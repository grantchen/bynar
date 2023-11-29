package config

// AppConfig is the interface for the application configuration
type AppConfig interface {
	GetDBConnection() string
	GetAccountManagementConnection() string
}
