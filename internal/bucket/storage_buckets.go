package bucket

import (
	"context"
	"sync"
	"time"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

type StorageBuckets struct {
	logger          common.LoggerInterface
	buckets         map[string]*TokenBucket
	capacity        float32
	refillRate      float32
	ttl             time.Duration
	cleanupInterval time.Duration
	mu              sync.RWMutex
	ctx             context.Context
}

func NewStorageBuckets(ctx *context.Context, cfg *config.LimitConfig, logger common.LoggerInterface) *StorageBuckets {
	storage := &StorageBuckets{
		logger:          logger,
		buckets:         make(map[string]*TokenBucket),
		capacity:        cfg.MaxPerMinute,
		refillRate:      cfg.RefillRateIsSecond,
		ttl:             cfg.TTL,
		cleanupInterval: cfg.CleanupInterval,
		mu:              sync.RWMutex{},
		ctx:             *ctx,
	}

	go storage.cleanup()

	return storage
}

func (s *StorageBuckets) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.buckets)
}

func (s *StorageBuckets) Allow(_ context.Context, key string) (bool, error) {
	s.mu.Lock()
	bucket, exists := s.buckets[key]
	if !exists {
		bucket = NewTokenBucket(s.capacity, s.refillRate)
		s.buckets[key] = bucket
	}
	s.mu.Unlock()

	return bucket.Request(1), nil
}

func (s *StorageBuckets) cleanup() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Start cleanup buckets")
			s.remove()
			s.logger.Info("Finish cleanup buckets")
		case <-s.ctx.Done():
			s.logger.Info("Stopping cleanup buckets")
			return
		}
	}
}

func (s *StorageBuckets) remove() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for key, bucket := range s.buckets {
		if bucket.isExpired(s.ttl) {
			delete(s.buckets, key)
		}
	}
	s.logger.Info("Сleanup buckets by TTL", "time", time.Since(now).Seconds(), "buckets", len(s.buckets))
}

func (s *StorageBuckets) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.buckets)
}
