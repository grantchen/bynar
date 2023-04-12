package logger

import (
	"fmt"
	"log"
)

var LogLevel = "Debug"

// Debug
func Debug(v ...interface{}) {
	if LogLevel != "Debug" {
		return
	}

	log.Output(2, fmt.Sprintln(v...))
}
