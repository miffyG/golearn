package logger

import (
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func Init() error {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	Sugar = Logger.Sugar()
	return nil
}

func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}
