package app

import (
	"context"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/bucket" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

type App struct {
	logger          common.LoggerInterface
	cfg             *config.Config
	ctx             *context.Context
	storageIP       *bucket.StorageBuckets
	storageLogin    *bucket.StorageBuckets
	storagePassword *bucket.StorageBuckets
}

func New(ctx *context.Context, cfg *config.Config, logger common.LoggerInterface) (*App, error) {
	storageLogin := bucket.NewStorageBuckets(ctx, &cfg.LimitsConfig.Login, logger)
	storagePassword := bucket.NewStorageBuckets(ctx, &cfg.LimitsConfig.Password, logger)
	storageIP := bucket.NewStorageBuckets(ctx, &cfg.LimitsConfig.IP, logger)
	logger.Debug("Config login", "cfg", cfg.LimitsConfig.Login)
	logger.Debug("Config password", "cfg", cfg.LimitsConfig.Password)
	logger.Debug("Config IP", "cfg", cfg.LimitsConfig.IP)
	return &App{
		cfg:             cfg,
		logger:          logger,
		ctx:             ctx,
		storageLogin:    storageLogin,
		storagePassword: storagePassword,
		storageIP:       storageIP,
	}, nil
}

func (app *App) CheckAuthLogin(login string) (bool, error) {
	return app.storageLogin.Allow(*app.ctx, login)
}

func (app *App) CheckAuthIP(ip string) (bool, error) {
	return app.storageIP.Allow(*app.ctx, ip)
}

func (app *App) CheckAuthPassword(password string) (bool, error) {
	return app.storagePassword.Allow(*app.ctx, password)
}
