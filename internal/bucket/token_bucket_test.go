package bucket

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTokenBucket(t *testing.T) {
	t.Parallel()

	t.Run("Capacity", func(t *testing.T) {
		tokenBucket := NewTokenBucket(2, 1)
		allowCount := 0
		for i := 0; i < 2; i++ {
			if tokenBucket.Request(1) {
				allowCount++
			}
		}
		require.Equal(t, 2, allowCount)
	})

	t.Run("Capacity overflow", func(t *testing.T) {
		tokenBucket := NewTokenBucket(2, 1)
		allowCount := 0
		for i := 0; i < 3; i++ {
			if tokenBucket.Request(1) {
				allowCount++
			}
		}

		require.Equal(t, 2, allowCount)
	})

	t.Run("RefillRate", func(t *testing.T) {
		tokenBucket := NewTokenBucket(1, 2)
		allowCount := 0
		for i := 0; i < 3; i++ {
			if tokenBucket.Request(1) {
				allowCount++
			}
			time.Sleep(500 * time.Millisecond)
		}
		require.Equal(t, 3, allowCount)
	})

	t.Run("RefillRate overflow", func(t *testing.T) {
		tokenBucket := NewTokenBucket(1, 1)
		allowCount := 0
		for i := 0; i < 3; i++ {
			if tokenBucket.Request(1) {
				allowCount++
			}
			time.Sleep(100 * time.Millisecond)
		}
		require.Equal(t, 1, allowCount)
	})
}
