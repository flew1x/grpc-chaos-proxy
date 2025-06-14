package injector

import (
	"context"
	"errors"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"testing"
	"time"
)

func TestNewDelay_InvalidConfig(t *testing.T) {
	_, err := NewDelay(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}

	_, err = NewDelay(123)
	if err == nil {
		t.Error("expected error for wrong type config")
	}
}

func TestNewDelay_MinMaxSwap(t *testing.T) {
	cfg := &config.DelayAction{MinMS: 100, MaxMS: 50}

	inj, err := NewDelay(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	di := inj.(*DelayInjector)
	if di.min > di.max {
		t.Error("min should not be greater than max after swap")
	}
}

func TestDelayInjector_Apply_NoDelay(t *testing.T) {
	cfg := &config.DelayAction{MinMS: 0, MaxMS: 0}
	inj, _ := NewDelay(cfg)
	frame := &engine.Frame{Ctx: context.Background()}
	start := time.Now()

	err := inj.Apply(frame)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if time.Since(start) > 10*time.Millisecond {
		t.Error("should not sleep if delay is zero")
	}
}

func TestDelayInjector_Apply_Delay(t *testing.T) {
	cfg := &config.DelayAction{MinMS: 10, MaxMS: 20}
	inj, _ := NewDelay(cfg)
	frame := &engine.Frame{Ctx: context.Background()}
	start := time.Now()

	err := inj.Apply(frame)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < 10*time.Millisecond || elapsed > 25*time.Millisecond {
		t.Errorf("delay out of expected range: %v", elapsed)
	}
}

func TestDelayInjector_Apply_ContextCancel(t *testing.T) {
	cfg := &config.DelayAction{MinMS: 100, MaxMS: 200}
	inj, _ := NewDelay(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	frame := &engine.Frame{Ctx: ctx}

	err := inj.Apply(frame)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context deadline exceeded, got: %v", err)
	}
}

func TestDelayInjector_Apply_WithDelay(t *testing.T) {
	cfg := &config.DelayAction{MinMS: 10, MaxMS: 20}
	inj, err := NewDelay(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := &engine.Frame{
		Ctx:     context.Background(),
		Service: "TestService",
		Method:  "TestMethod",
		MD:      nil,
	}

	start := time.Now()
	err = inj.Apply(frame)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if elapsed < 10*time.Millisecond {
		t.Errorf("expected at least 10ms delay, got %v", elapsed)
	}
}
