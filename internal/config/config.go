package config

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/heetch/confita"              //nolint:depguard
	"github.com/heetch/confita/backend/env"  //nolint:depguard
	"github.com/heetch/confita/backend/file" //nolint:depguard
)

var (
	ErrLoggerLevel = errors.New("logger Level")
	ErrLoadConfig  = errors.New("loading config")
)

type Config struct {
	App          AppConfig
	Logger       LogConfig
	Storage      StorageConfig
	HTTP         HTTPConfig
	LimitsConfig LimitsConfig `yaml:"limits" env:"LIMITS"`
}

type AppConfig struct {
	Overlapped bool `yaml:"overlapped" env:"OVERLAPPED"`
}

type LogConfig struct {
	Level string `config:"level"`
}

type LimitsConfig struct {
	Login    LimitConfig `yaml:"login" env-prefix:"LOGIN_"`
	Password LimitConfig `yaml:"password" env-prefix:"PASSWORD_"`
	IP       LimitConfig `yaml:"ip" env-prefix:"IP_"`
}

type LimitConfig struct {
	MaxPerMinute       float32       `yaml:"maxPerMinute" env:"MAX_PER_MINUTE"`
	RefillRateIsSecond float32       `yaml:"refillRateInSecond" env:"REFILL_RATE_IN_SECOND"`
	CleanupInterval    time.Duration `yaml:"cleanupInterval" env:"CLEANUP_INTERVAL"`
	TTL                time.Duration `yaml:"ttl" env:"TTL"`
}

type StorageConfig struct {
	Type string `config:"type"`
	Dsn  string `config:"dsn"`
}

type HTTPConfig struct {
	Host         string        `config:"host"`
	Port         string        `config:"port"`
	ReadTimeout  time.Duration `config:"readTimeout"`
	WriteTimeout time.Duration `config:"writeTimeout"`
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

func New(configPath string) (*Config, error) {
	loggerLeverPosible := []string{"info", "warn", "debug", "error", ""}
	cfg := Config{
		App: AppConfig{},
		Logger: LogConfig{
			Level: "info",
		},
		Storage: StorageConfig{},
		HTTP: HTTPConfig{
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		LimitsConfig: LimitsConfig{},
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
		return &cfg, err
	}

	if !slices.Contains(loggerLeverPosible, cfg.Logger.Level) {
		return &cfg, fmt.Errorf("config load: %w", ErrLoggerLevel)
	}

	return &cfg, err
}
