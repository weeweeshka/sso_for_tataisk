package main

import (
	"github.com/weeweeshka/sso_for_tataisk/internal/config"
	"github.com/weeweeshka/sso_for_tataisk/pkg/libs/logger"
)

func main() {
	cfg := config.MustLoad()
	logr := logger.SetupLogger()

	logr.Info("Config loaded")
	logr.Info("Logger initialized")

}
