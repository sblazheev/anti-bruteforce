package app

import (
	"context"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

type App struct {
	dateOverlapping bool
	logger          common.LoggerInterface
	cfg             *config.Config
	ctx             *context.Context
}

func New(cfg *config.Config, logger common.LoggerInterface, ctx *context.Context) (*App, error) {
	return &App{
		dateOverlapping: cfg.App.Overlapping,
		cfg:             cfg,
		logger:          logger,
		ctx:             ctx,
	}, nil
}
