package buildApp

import (
	runGrpc "github.com/weeweeshka/sso_for_tataisk/internal/app/grpcApp"
	"github.com/weeweeshka/sso_for_tataisk/internal/repository"
	"github.com/weeweeshka/sso_for_tataisk/internal/services/sso"
	"go.uber.org/zap"
	"time"
)

type App struct {
	Grpc runGrpc.GRPCServer
}

func NewApp(port int, storagePath string, logr *zap.Logger, tokenTTL time.Duration) (*App, error) {

	storage, err := repository.NewStorage(storagePath, logr)
	if err != nil {
		return &App{}, err
	}

	ssoService := sso.NewSsoService(logr, storage, tokenTTL)
	grpc := runGrpc.New(port, logr, ssoService)

	return &App{Grpc: *grpc}, nil
}
