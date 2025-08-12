package logger

import "go.uber.org/zap"

func SetupLogger() *zap.Logger {
	logr, err := zap.NewDevelopment()
	if err != nil {
		panic("impossible to setup logger" + err.Error())
	}

	return logr
}
