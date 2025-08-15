package main

import (
	"github.com/weeweeshka/sso_for_tataisk/internal/app/buildApp"
	"github.com/weeweeshka/sso_for_tataisk/internal/config"
	"github.com/weeweeshka/sso_for_tataisk/pkg/libs/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	logr := logger.SetupLogger()

	logr.Info("Config loaded")
	logr.Info("Logger initialized")

	app, err := buildApp.NewApp(cfg.GRPC.Port, cfg.StoragePath, logr, cfg.TokenTTL)
	if err != nil {
		logr.Info("Error initializing app")
		panic(err)
	}

	go func() {
		if err := app.Grpc.Run(); err != nil {
			logr.Info("failed to run", zap.Error(err))
		}
	}()
	logr.Info("Server started", zap.Int("port", cfg.GRPC.Port))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	logr.Info("App stopped")
	app.Grpc.GracefulStop()
}
