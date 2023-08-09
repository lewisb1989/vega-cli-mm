package logging

import (
	"github.com/withmandala/go-log"
	"os"
)

func GetLogger() *log.Logger {
	if os.Getenv("DEBUG") == "1" {
		return log.New(os.Stderr).WithDebug().WithColor()
	} else {
		return log.New(os.Stderr).WithColor()
	}
}

func Panic(message string) {
	GetLogger().Errorf(message)
	panic(message)
}
