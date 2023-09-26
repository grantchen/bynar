/**
    @author: dongjs
    @date: 2023/9/14
    @description:
**/

package config

import "os"

// IsProductionEnv Whether the code is running in a production environment
// this is default
func IsProductionEnv() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == ""
}

// IsDebugEnv Whether the code is running in a debug environment
func IsDebugEnv() bool {
	return os.Getenv("ENV") == "debug"
}

// IsDevelopmentEnv Whether the code is running in a development environment
func IsDevelopmentEnv() bool {
	return os.Getenv("ENV") == "development"
}
