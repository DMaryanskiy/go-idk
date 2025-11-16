package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() *zap.Logger {
	var logger *zap.Logger
	var err error

	if os.Getenv("ENV") == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	return logger
}
