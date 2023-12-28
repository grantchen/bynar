package logger

import (
	"fmt"
	"log"
)

func Debug(v ...interface{}) {
	_ = log.Output(2, fmt.Sprintln(v...))
}
