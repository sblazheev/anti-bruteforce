package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/logger" //nolint:depguard
)

func TestApp(t *testing.T) {
	c, err := config.New("./test/config.yaml")
	require.NoError(t, err)
	logg := logger.New(&c.Logger)

	ctx := context.Background()

	app, err := New(c, logg, &ctx)
	require.NoError(t, err)

	t.Run("Add event", func(t *testing.T) {
		require.NotNil(t, app)
		require.Equal(t, true, true)
	})
}
