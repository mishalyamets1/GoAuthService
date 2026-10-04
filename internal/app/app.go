package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"

)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpsPort int,
	storagePath string, 
	tokenTTL time.Duration,
) (*App) {
	grpcApp := grpcapp.New(log, grpsPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}


