package network

import (
	"context"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"testing"
	"time"
)

func TestNetworkInjector_Apply_NoLossNoThrottle(t *testing.T) {
	inj := &Injector{LossPercentage: 0, ThrottleMS: 0}
	frame := &engine.Frame{Ctx: context.Background()}

	if err := inj.Apply(frame); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestNetworkInjector_Apply_LossAlways(t *testing.T) {
	inj := &Injector{LossPercentage: 100, ThrottleMS: 0}
	frame := &engine.Frame{Ctx: context.Background()}

	for i := 0; i < 10; i++ {
		err := inj.Apply(frame)
		if err == nil {
			t.Error("expected error due to 100% loss, got nil")
		}
	}
}

func TestNetworkInjector_Apply_LossNever(t *testing.T) {
	inj := &Injector{LossPercentage: 0, ThrottleMS: 0}
	frame := &engine.Frame{Ctx: context.Background()}

	for i := 0; i < 10; i++ {
		err := inj.Apply(frame)
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	}
}

func TestNetworkInjector_Apply_Throttle(t *testing.T) {
	inj := &Injector{LossPercentage: 0, ThrottleMS: 50}
	frame := &engine.Frame{Ctx: context.Background()}

	start := nowMillis()
	err := inj.Apply(frame)
	elapsed := nowMillis() - start

	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}

	if elapsed < 45 {
		t.Errorf("expected at least 45ms throttle, got %dms", elapsed)
	}
}

func nowMillis() int64 {
	return int64((1e-6) * float64(time.Now().UnixNano()))
}
