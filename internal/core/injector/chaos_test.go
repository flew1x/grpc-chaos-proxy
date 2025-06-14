package injector

import (
	"context"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"testing"
)

func TestChaosInjector_Apply_OneAction(t *testing.T) {
	delay := &config.DelayAction{MinMS: 1, MaxMS: 1}
	chaosCfg := &config.ChaosAction{
		Actions: []config.Action{{Delay: delay}},
	}

	inj, err := NewChaosInjector(chaosCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := &engine.Frame{
		Ctx:     context.Background(),
		Service: "TestService",
		Method:  "TestMethod",
		MD:      nil,
	}

	if err := inj.Apply(frame); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestChaosInjector_Apply_MultipleActions(t *testing.T) {
	delay := &config.DelayAction{MinMS: 1, MaxMS: 1}
	abort := &config.AbortAction{Code: "internal", Percentage: 100}

	chaosCfg := &config.ChaosAction{
		Actions: []config.Action{
			{Delay: delay},
			{Abort: abort},
		},
	}

	inj, err := NewChaosInjector(chaosCfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := &engine.Frame{Ctx: context.Background()}

	for i := 0; i < 10; i++ {
		_ = inj.Apply(frame)
	}
}

func TestChaosInjector_EmptyConfig(t *testing.T) {
	chaosCfg := &config.ChaosAction{Actions: nil}
	_, err := NewChaosInjector(chaosCfg)
	if err == nil {
		t.Error("expected error for empty actions, got nil")
	}
}

func TestChaosInjector_InvalidType(t *testing.T) {
	_, err := NewChaosInjector(nil)
	if err == nil {
		t.Error("expected error for nil config, got nil")
	}
}
