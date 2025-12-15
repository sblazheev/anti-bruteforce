package config

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/heetch/confita"              //nolint:depguard
	"github.com/heetch/confita/backend/env"  //nolint:depguard
	"github.com/heetch/confita/backend/file" //nolint:depguard
)

var (
	ErrLoggerLevel = errors.New("logger Level")
	ErrLoadConfig  = errors.New("loading config")
)

type Config struct {
	App       AppConfig
	Logger    LogConfig
	Storage   StorageConfig
	HTTP      HTTPConfig
	Grpc      GrpcConfig
	Broker    BrokerConfig
	Scheduler SchedulerConfig
	Sender    SenderConfig
}

type AppConfig struct {
	Overlapping bool `config:"overlapping"`
}

type LogConfig struct {
	Level string `config:"level"`
}

type StorageConfig struct {
	Type string `config:"type"`
	Dsn  string `config:"dsn"`
}

type HTTPConfig struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

type GrpcConfig struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

type BrokerConfig struct {
	Queue QueueConfig
	Ampq  string `config:"ampq"`
}

type QueueConfig struct {
	Notify string `config:"notify"`
}

type SchedulerConfig struct {
	Interval int `config:"interval"`
	KeepDays int `config:"keepdays"`
	Chunk    int `config:"chunk"`
}

type SenderConfig struct {
	Interval int `config:"interval"`
}

func New(configPath string) (*Config, error) {
	loggerLeverPosible := []string{"info", "warn", "debug", "error", ""}
	cfg := Config{
		App: AppConfig{
			Overlapping: true,
		},
		Logger: LogConfig{
			Level: "info",
		},
		Storage:   StorageConfig{},
		HTTP:      HTTPConfig{},
		Grpc:      GrpcConfig{},
		Broker:    BrokerConfig{},
		Scheduler: SchedulerConfig{Chunk: 100, KeepDays: 365, Interval: 10},
		Sender:    SenderConfig{Interval: 10},
	}
	var loader *confita.Loader
	if len(configPath) > 0 {
		loader = confita.NewLoader(
			file.NewBackend(configPath),
			env.NewBackend(),
		)
	} else {
		loader = confita.NewLoader(
			env.NewBackend(),
		)
	}

	err := loader.Load(context.Background(), &cfg)
	if err != nil {
		return &cfg, ErrLoadConfig
	}

	if !slices.Contains(loggerLeverPosible, cfg.Logger.Level) {
		return &cfg, fmt.Errorf("config load: %w", ErrLoggerLevel)
	}

	return &cfg, err
}
