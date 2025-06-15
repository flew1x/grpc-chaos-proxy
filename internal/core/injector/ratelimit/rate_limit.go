package ratelimit

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"sync"
	"time"
)

type RateLimit struct {
	Limit int
	Burst int

	mu        sync.Mutex
	allowance int
	lastCheck time.Time
}

func NewRateLimit(cfg any) (engine.Injector, error) {
	ac, ok := cfg.(*config.RateLimiterAction)
	if !ok {
		return nil, apperr.ErrInvalidConfig
	}

	if ac.RateLimit <= 0 {
		return nil, fmt.Errorf("ratelimit action config error: RateLimit must be greater than 0")
	}

	if ac.BurstSize < 0 {
		ac.BurstSize = 0
	}

	return &RateLimit{
		Limit: ac.RateLimit,
		Burst: ac.BurstSize,
	}, nil

}

func (rl *RateLimit) Apply(frame *engine.Frame) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if rl.lastCheck.IsZero() {
		rl.lastCheck = now
		rl.allowance = rl.Limit + rl.Burst
	}

	secElapsed := now.Sub(rl.lastCheck).Seconds()

	if secElapsed > 0 {
		newTokens := int(secElapsed * float64(rl.Limit))
		rl.allowance += newTokens

		if rl.allowance > rl.Limit+rl.Burst {
			rl.allowance = rl.Limit + rl.Burst
		}

		rl.lastCheck = now
	}

	if rl.allowance > 0 {
		rl.allowance--

		return nil
	}

	return apperr.ErrRateLimitExceeded
}

func init() {
	engine.Register(entity.RateLimitType, NewRateLimit)
}
