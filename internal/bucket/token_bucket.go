package bucket

import (
	"sync"
	"time"
)

// TokenBucket represents a token bucket system.
type TokenBucket struct {
	lastRefillTime time.Time
	mu             sync.RWMutex
	tokens         float32
	capacity       float32
	refillRate     float32
}

func NewTokenBucket(capacity float32, refillRateInSecond float32) *TokenBucket {
	return &TokenBucket{
		tokens:         capacity,
		capacity:       capacity,
		refillRate:     refillRateInSecond,
		lastRefillTime: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	duration := now.Sub(tb.lastRefillTime)
	tokens := tb.tokens + tb.refillRate*float32(duration.Seconds())
	if tokens > tb.capacity {
		tb.tokens = tb.capacity
	} else {
		tb.tokens = tokens
	}
	tb.lastRefillTime = now
}

func (tb *TokenBucket) lastRefill() time.Time {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	return tb.lastRefillTime
}

func (tb *TokenBucket) isExpired(ttl time.Duration) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	now := time.Now()
	if now.Sub(tb.lastRefill()) > ttl {
		return true
	}
	return false
}

func (tb *TokenBucket) Request(tokens float32) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	if tokens <= tb.tokens {
		tb.tokens -= tokens
		return true
	}
	return false
}
