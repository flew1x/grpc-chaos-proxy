package disconnect

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestDisconnect_ZeroPercentage(t *testing.T) {
	inj := Disconnect{Percentage: 0}

	err := inj.Apply(&engine.Frame{})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestDisconnect_AlwaysDisconnect(t *testing.T) {
	inj := Disconnect{Percentage: 100}

	for i := 0; i < 10; i++ {
		err := inj.Apply(&engine.Frame{})
		if err == nil {
			t.Error("expected error, got nil")
		}

		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Unavailable {
			t.Errorf("expected UNAVAILABLE code, got %v", err)
		}
	}
}

func TestDisconnect_NeverDisconnect(t *testing.T) {
	inj := Disconnect{Percentage: 0}

	for i := 0; i < 10; i++ {
		err := inj.Apply(&engine.Frame{})
		if err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	}
}

func TestDisconnect_RandomDisconnect(t *testing.T) {
	inj := Disconnect{Percentage: 50}
	aborted := 0
	total := 1000

	for i := 0; i < total; i++ {
		err := inj.Apply(&engine.Frame{})
		if err != nil {
			aborted++
		}
	}

	if aborted < 400 || aborted > 600 {
		t.Errorf("expected ~50%% disconnects, got %d", aborted)
	}
}
