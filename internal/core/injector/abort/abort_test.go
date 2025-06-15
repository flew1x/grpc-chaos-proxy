package abort

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"testing"
)

func TestNewAbort_InvalidConfig(t *testing.T) {
	_, err := NewAbortInjector(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}

	_, err = NewAbortInjector(123)
	if err == nil {
		t.Error("expected error for wrong type config")
	}
}

func TestAbortInjector_Apply_PercentageZero(t *testing.T) {
	cfg := &config.AbortAction{Code: "UNAVAILABLE", Percentage: 0}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = inj.Apply(&engine.Frame{})
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestAbortInjector_Apply_AlwaysAbort(t *testing.T) {
	cfg := &config.AbortAction{Code: "UNAVAILABLE", Percentage: 100}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 10; i++ {
		err = inj.Apply(&engine.Frame{})
		if err == nil {
			t.Error("expected error, got nil")
		}

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Unavailable {
			t.Errorf("expected UNAVAILABLE code, got: %v", err)
		}
	}
}

func TestAbortInjector_Apply_NeverAbort(t *testing.T) {
	cfg := &config.AbortAction{Code: "UNAVAILABLE", Percentage: 0}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 10; i++ {
		err = inj.Apply(&engine.Frame{})
		if err != nil {
			t.Errorf("expected nil error, got: %v", err)
		}
	}
}

func TestAbortInjector_Apply_RandomAbort(t *testing.T) {
	cfg := &config.AbortAction{Code: "UNAVAILABLE", Percentage: 50}

	inj, err := NewAbortInjector(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	aborted := 0
	total := 1000
	rand.New(rand.NewSource(64))

	for i := 0; i < total; i++ {
		err = inj.Apply(&engine.Frame{})
		if err != nil {
			aborted++
		}
	}

	if aborted < 400 || aborted > 600 {
		t.Errorf("expected ~50%% aborts, got %d", aborted)
	}
}
