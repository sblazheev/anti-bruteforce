package logger

import (
	"testing"

	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

func TestLogger(t *testing.T) {
	t.Run("Logger create", func(t *testing.T) {
		logger := New(&config.LogConfig{
			Level: "info",
		})
		require.Implements(t, (*common.LoggerInterface)(nil), logger)
	})
}
