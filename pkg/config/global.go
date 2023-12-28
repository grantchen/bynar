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
