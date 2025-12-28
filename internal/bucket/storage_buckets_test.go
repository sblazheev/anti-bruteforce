package bucket

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"                     //nolint:depguard
	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/logger" //nolint:depguard
)

func printMemoryUsage(numberBuckets int) string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return fmt.Sprintf("Total Buckets = %v Pcs", numberBuckets) +
		fmt.Sprintf("\tAlloc = %v MiB", m.Alloc/1024/1024) +
		fmt.Sprintf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024) +
		fmt.Sprintf("\tSys = %v MiB", m.Sys/1024/1024) +
		fmt.Sprintf("\tNumGC = %v\n", m.NumGC)
}

func TestStorageBuckets(t *testing.T) {
	t.Parallel()

	capacity := 1_00
	CleanupInterval := 100 * time.Millisecond
	TTL := 1 * time.Microsecond

	logins := make(map[string]struct{}, capacity)
	for len(logins) < capacity {
		logins[gofakeit.Email()] = struct{}{}
	}

	cfg, _ := config.New("")
	cfg.Logger.Level = "error"
	log := logger.New(&cfg.Logger)

	t.Run("Ban", func(t *testing.T) {
		t.Parallel()

		allowCount := 0
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		storageBuckets := NewStorageBuckets(
			&ctx,
			"test",
			&config.LimitConfig{MaxPerMinute: 1, RefillRateIsSecond: 0.0000001, CleanupInterval: CleanupInterval, TTL: TTL},
			log,
		)

		for login := range logins {
			allow, _ := storageBuckets.Allow(ctx, login)
			if allow {
				allowCount++
			}
		}
		require.Equal(t, len(logins), allowCount)
		require.Equal(t, len(logins), storageBuckets.Len())

		allowCount = 0
		for login := range logins {
			allow, _ := storageBuckets.Allow(ctx, login)
			if allow {
				allowCount++
			}
		}
		cancel()
		storageBuckets.reset()
		require.Equal(t, 0, allowCount)
		log.Info("Memory", "stats", printMemoryUsage(len(logins)))
	})

	t.Run("CleanUp", func(t *testing.T) {
		t.Parallel()

		allowCount := 0
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		storageBuckets := NewStorageBuckets(
			&ctx,
			"test",
			&config.LimitConfig{MaxPerMinute: 1, RefillRateIsSecond: 1, CleanupInterval: CleanupInterval, TTL: TTL},
			log,
		)
		for i := 0; i < 10; i++ {
			allow, _ := storageBuckets.Allow(ctx, gofakeit.Email())
			if allow {
				allowCount++
			}
		}

		require.Equal(t, 10, storageBuckets.Len())

		time.Sleep(CleanupInterval * 2)
		require.Equal(t, 0, storageBuckets.Len())
		cancel()
		log.Info("Memory", "stats", printMemoryUsage(len(logins)))
	})
}
