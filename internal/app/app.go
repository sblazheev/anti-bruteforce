package app

import (
	"context"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/bucket"                           //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"                           //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config"                           //nolint:depguard
	memorystorage "gitlab.wsrubi.ru/go/anti-bruteforce/internal/storage/memory"     //nolint:depguard
	sqlstorage "gitlab.wsrubi.ru/go/anti-bruteforce/internal/storage/sql"           //nolint:depguard
	sqlproxystorage "gitlab.wsrubi.ru/go/anti-bruteforce/internal/storage/sqlproxy" //nolint:depguard
)

type App struct {
	logger           common.LoggerInterface
	cfg              *config.Config
	ctx              *context.Context
	storageWhiteList *common.Storage
	storageBlackList *common.Storage
	storageIP        *bucket.StorageBuckets
	storageLogin     *bucket.StorageBuckets
	storagePassword  *bucket.StorageBuckets
}

func NewStorageDriver(ctx *context.Context, c config.StorageConfig) (common.StorageDriverInterface, error) {
	switch c.Type {
	case "sqlproxy":
		return sqlproxystorage.New(ctx, c), nil
	case "sql":
		return sqlstorage.New(ctx, c), nil
	case "memory":
		return memorystorage.New(), nil
	default:
		return memorystorage.New(), nil
	}
}

func New(ctx *context.Context, cfg *config.Config, logger common.LoggerInterface) (*App, error) {
	storageDriverWhite, err := NewStorageDriver(ctx, cfg.Storage)
	if err != nil {
		return nil, err
	}

	storageDriverBlack, err := NewStorageDriver(ctx, cfg.Storage)
	if err != nil {
		return nil, err
	}

	storageWhiteList, err := common.NewStorage("white_list", ctx, storageDriverWhite)
	if err != nil {
		return nil, err
	}
	if _, err := storageWhiteList.Load(); err != nil {
		return nil, err
	}

	storageBlackList, err := common.NewStorage("black_list", ctx, storageDriverBlack)
	if err != nil {
		return nil, err
	}
	if _, err := storageBlackList.Load(); err != nil {
		return nil, err
	}

	storageLogin := bucket.NewStorageBuckets(ctx, "Login", &cfg.LimitsConfig.Login, logger)
	storagePassword := bucket.NewStorageBuckets(ctx, "Password", &cfg.LimitsConfig.Password, logger)
	storageIP := bucket.NewStorageBuckets(ctx, "IP", &cfg.LimitsConfig.IP, logger)
	logger.Debug("Config login", "cfg", cfg.LimitsConfig.Login)
	logger.Debug("Config password", "cfg", cfg.LimitsConfig.Password)
	logger.Debug("Config IP", "cfg", cfg.LimitsConfig.IP)
	return &App{
		cfg:              cfg,
		logger:           logger,
		ctx:              ctx,
		storageLogin:     storageLogin,
		storagePassword:  storagePassword,
		storageIP:        storageIP,
		storageWhiteList: storageWhiteList,
		storageBlackList: storageBlackList,
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

func (app *App) CheckWhiteList(ip string) (bool, error) {
	return app.storageWhiteList.InSubNet(ip)
}

func (app *App) CheckBlackList(ip string) (bool, error) {
	return app.storageBlackList.InSubNet(ip)
}

func (app *App) DeleteIPBucket(ip string) bool {
	app.storageIP.RemoveBucket(ip)
	return true
}

func (app *App) DeleteLoginBucket(ip string) bool {
	app.storageLogin.RemoveBucket(ip)
	return true
}

func (app *App) AddBlackList(net string) (*common.IPSubnet, error) {
	return app.storageBlackList.Add(common.IPSubnet{
		Subnet: net,
	})
}

func (app *App) AddWhiteList(net string) (*common.IPSubnet, error) {
	return app.storageWhiteList.Add(common.IPSubnet{
		Subnet: net,
	})
}

func (app *App) DeleteBlackList(net string) error {
	return app.storageBlackList.Delete(net)
}

func (app *App) DeleteWhiteList(net string) error {
	return app.storageWhiteList.Delete(net)
}
