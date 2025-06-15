package ratelimit

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"testing"
	"time"
)

func newRateLimitInjector(limit, burst int) *RateLimit {
	rl := &RateLimit{Limit: limit, Burst: burst}
	rl.allowance = limit + burst
	rl.lastCheck = time.Now()

	return rl
}

func TestRateLimit_AllowWithinLimit(t *testing.T) {
	rl := newRateLimitInjector(2, 1)
	frame := &engine.Frame{}

	if err := rl.Apply(frame); err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if err := rl.Apply(frame); err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if err := rl.Apply(frame); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRateLimit_ExceedLimit(t *testing.T) {
	rl := newRateLimitInjector(1, 1)
	frame := &engine.Frame{}

	_ = rl.Apply(frame)
	_ = rl.Apply(frame)

	err := rl.Apply(frame)
	if err == nil {
		t.Error("expected error when exceeding rate limit")
	}
}

func TestRateLimit_RefillTokens(t *testing.T) {
	rl := newRateLimitInjector(1, 0)
	frame := &engine.Frame{}
	_ = rl.Apply(frame)

	err := rl.Apply(frame)
	if err == nil {
		t.Error("expected error when exceeding rate limit")
	}

	time.Sleep(1100 * time.Millisecond)

	if err := rl.Apply(frame); err != nil {
		t.Errorf("expected nil after refill, got %v", err)
	}
}

func TestRateLimit_InvalidConfig(t *testing.T) {
	_, err := NewRateLimit(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}

	_, err = NewRateLimit(123)
	if err == nil {
		t.Error("expected error for wrong type config")
	}

	_, err = NewRateLimit(&config.RateLimiterAction{RateLimit: 0})
	if err == nil {
		t.Error("expected error for zero rate limit")
	}
}
