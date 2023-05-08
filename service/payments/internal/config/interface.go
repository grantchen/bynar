package config

type AppConfig interface {
	GetDBConnection() string
}
