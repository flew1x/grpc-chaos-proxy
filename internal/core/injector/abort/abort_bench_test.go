package abort

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"testing"
)

func BenchmarkAbortInjector_Apply(b *testing.B) {
	cfg := &config.AbortAction{
		Code:       "internal",
		Percentage: 100,
	}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		b.Fatalf("failed to create abort injector: %v", err)
	}

	frame := &engine.Frame{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = inj.Apply(frame)
	}
}

func BenchmarkAbortInjector_Alloc(b *testing.B) {
	cfg := &config.AbortAction{
		Code:       "internal",
		Percentage: 100,
	}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		b.Fatalf("failed to create abort injector: %v", err)
	}

	frame := &engine.Frame{}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = inj.Apply(frame)
	}
}

func TestAbortInjector_Apply(t *testing.T) {
	cfg := &config.AbortAction{
		Code:       "internal",
		Percentage: 100,
	}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		t.Fatalf("failed to create abort injector: %v", err)
	}

	frame := &engine.Frame{}
	err = inj.Apply(frame)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if _, ok := err.(interface{ GRPCStatus() interface{} }); !ok && err.Error() == "" {
		t.Errorf("unexpected error type: %T", err)
	}
}
