package middlewares

import (
	"api-versioning-me/config"
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func LoggerInit() (err error) {
	if config.Env == "production" {
		Logger, err = zapdriver.NewProduction()
	} else {
		Logger, err = zapdriver.NewDevelopment()
	}
	return err
}

func GetLogger() *zap.Logger {
	return Logger
}
