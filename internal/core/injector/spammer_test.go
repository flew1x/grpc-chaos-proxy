package injector

import (
	"context"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"

	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
)

func TestSpammerInjector_Apply_NoDelay(t *testing.T) {
	action := &config.SpammerAction{Count: 3}

	calls := 0
	inj, err := NewSpammerInjector(action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	injTyped := inj.(*SpammerInjector)
	injTyped.sender = func(ctx context.Context, service, method string, md metadata.MD) error {
		calls++

		return nil
	}

	frame := &engine.Frame{Ctx: context.Background(), Service: "svc", Method: "m"}

	if err := inj.Apply(frame); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestSpammerInjector_Apply_WithDelay(t *testing.T) {
	action := &config.SpammerAction{Count: 2, DelayAction: &config.DelayAction{MinMS: 10, MaxMS: 20}}

	inj, err := NewSpammerInjector(action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	calls := 0
	injTyped := inj.(*SpammerInjector)

	injTyped.sender = func(ctx context.Context, service, method string, md metadata.MD) error {
		calls++

		return nil
	}

	frame := &engine.Frame{Ctx: context.Background(), Service: "svc", Method: "m"}
	start := time.Now()

	if err := inj.Apply(frame); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	elapsed := time.Since(start)

	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}

	if elapsed < 10*time.Millisecond {
		t.Errorf("expected at least 10ms delay, got %v", elapsed)
	}
}

func TestNewSpammerInjector_InvalidConfig(t *testing.T) {
	_, err := NewSpammerInjector(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}

	_, err = NewSpammerInjector(&config.SpammerAction{Count: 0})
	if err == nil {
		t.Error("expected error for zero count")
	}
}

func TestSpammerInjector_Apply_Integration(t *testing.T) {
	calls := 0
	action := &config.SpammerAction{Count: 5, DelayAction: &config.DelayAction{MinMS: 0, MaxMS: 0}}

	inj, err := NewSpammerInjector(action)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	injTyped := inj.(*SpammerInjector)

	injTyped.sender = func(ctx context.Context, service, method string, md metadata.MD) error {
		calls++
		if service != "svc" || method != "m" {
			t.Errorf("unexpected service/method: %s/%s", service, method)
		}

		return nil
	}

	frame := &engine.Frame{Ctx: context.Background(), Service: "svc", Method: "m"}
	if err := inj.Apply(frame); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if calls != 5 {
		t.Errorf("expected 5 calls, got %d", calls)
	}
}
