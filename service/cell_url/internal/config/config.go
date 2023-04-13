package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	mySqlHost       string
	mySqlPort       string
	mySqlDBName     string
	mySqlDBPassword string
	mySqlUsername   string
}

const (
	ENVMySqlHost       = "mysql.dbhost"
	ENVMySqlPort       = "mysql.dbport"
	ENVMySqlDBName     = "mysql.dbname"
	ENVMySqlDBPassword = "mysql.dbpassword"
	ENVMySqlUsername   = "mysql.username"
)

var appConfig config

func Load() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
		return err
	}
	appConfig = config{
		mySqlHost:       os.Getenv(ENVMySqlHost),
		mySqlDBPassword: os.Getenv(ENVMySqlDBPassword),
		mySqlDBName:     os.Getenv(ENVMySqlDBName),
		mySqlUsername:   os.Getenv(ENVMySqlUsername),
		mySqlPort:       os.Getenv(ENVMySqlPort),
	}
	return nil
}

func GetMySqlHost() string {
	return appConfig.mySqlHost
}

func GetMySqlPort() string {
	return appConfig.mySqlPort
}

func GetMySqlDBName() string {
	return appConfig.mySqlDBName
}

func GetMySqlUsername() string {
	return appConfig.mySqlUsername
}

func GetMySqlPassword() string {
	return appConfig.mySqlDBPassword
}
