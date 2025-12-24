package config

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestConfig(t *testing.T) {
	t.Run("Config create", func(t *testing.T) {
		config, err := New("./test/config.yaml")
		require.Equal(t, &Config{
			App: AppConfig{},
			Logger: LogConfig{
				Level: "info",
			},
		}, config)
		require.NoError(t, err)
	})
	t.Run("Config logger Level error", func(t *testing.T) {
		config, err := New("./test/config2e.yaml")
		require.Equal(t, "info2", config.Logger.Level)
		require.ErrorIs(t, err, ErrLoggerLevel)
	})
}
