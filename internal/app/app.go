package app

import (
	"context"
	"fmt"
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/config"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
	storage *postgres.Storage
}

func New(
	log *slog.Logger,
	grpsPort int,
	pgCfg config.PostgresConfig, 
	tokenTTL time.Duration,
) (*App) {

	const op = "app.New"

	storage, err := postgres.New(
		context.Background(),
		pgCfg.Host,
		pgCfg.Port,
		pgCfg.User,
		pgCfg.Password,
		pgCfg.DB,
		pgCfg.Timeout,
	)
	if err != nil {
		panic(fmt.Errorf("%s: %w", op, err))
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, grpsPort, authService)
	return &App{
		GRPCSrv: grpcApp,
		storage: storage,
	}
}
func (a *App) Stop() {
	a.GRPCSrv.Stop()
	a.storage.Close()
}


