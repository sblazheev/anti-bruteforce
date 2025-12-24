package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/logger" //nolint:depguard
)

func TestApp(t *testing.T) {
	c, err := config.New("")
	c.LimitsConfig = config.LimitsConfig{Login: config.LimitConfig{
		MaxPerMinute:       1,
		CleanupInterval:    2 * 60 * time.Second,
		RefillRateIsSecond: 0.00001,
		TTL:                2 * 60 * time.Second,
	}}

	logg := logger.New(&config.LogConfig{})

	ctx := context.Background()

	app, err := New(&ctx, c, logg)
	require.NoError(t, err)

	t.Run("Allow Login", func(t *testing.T) {
		allow, err := app.CheckAuthLogin("test@test.ru")
		require.NoError(t, err)
		require.Equal(t, true, allow)
		allow, err = app.CheckAuthLogin("test@test.ru")
		require.NoError(t, err)
		require.Equal(t, false, allow)
	})
}
