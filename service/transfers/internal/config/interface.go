package config

type AppConfig interface {
	GetDBConnection() string
}

var (
	ModuleID = 6
)
