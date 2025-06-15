package code

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func newFrame(method string) *engine.Frame {
	return &engine.Frame{
		Method: method,
		MD:     metadata.New(nil),
	}
}

func TestCodeInjector_Percentage(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", Percentage: 100}
	frame := newFrame("TestMethod")

	err := inj.Apply(frame)
	if err == nil {
		t.Error("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Unavailable {
		t.Errorf("expected UNAVAILABLE, got %v", err)
	}
}

func TestCodeInjector_PercentageZero(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", Percentage: 0}
	frame := newFrame("TestMethod")

	err := inj.Apply(frame)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCodeInjector_OnlyOn(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", OnlyOn: []string{"Foo"}}
	frame := newFrame("Bar")

	err := inj.Apply(frame)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	frame = newFrame("Foo")

	err = inj.Apply(frame)
	if err == nil {
		t.Error("expected error for method Foo")
	}
}

func TestCodeInjector_Delay(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", DelayMS: 50}
	frame := newFrame("TestMethod")
	start := time.Now()
	_ = inj.Apply(frame)

	elapsed := time.Since(start)
	if elapsed < 45*time.Millisecond {
		t.Errorf("expected at least 45ms delay, got %v", elapsed)
	}
}

func TestCodeInjector_Metadata(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", Metadata: map[string]string{"x-test": "abc"}}
	frame := newFrame("TestMethod")

	_ = inj.Apply(frame)
	if got := frame.MD.Get("x-test"); len(got) == 0 || got[0] != "abc" {
		t.Errorf("expected metadata x-test=abc, got %v", got)
	}
}

func TestCodeInjector_RepeatCount(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", RepeatCount: 2}
	frame := newFrame("TestMethod")

	for i := 0; i < 2; i++ {
		err := inj.Apply(frame)
		if err == nil {
			t.Errorf("expected error on repeat %d", i)
		}
	}

	err := inj.Apply(frame)
	if err != nil {
		t.Errorf("expected nil after repeat count exceeded, got %v", err)
	}
}

func TestCodeInjector_CustomMessage(t *testing.T) {
	inj := &CodeInjector{Code: "UNAVAILABLE", Message: "custom msg"}
	frame := newFrame("TestMethod")

	err := inj.Apply(frame)
	st, _ := status.FromError(err)

	if st.Message() != "custom msg" {
		t.Errorf("expected custom message, got %v", st.Message())
	}
}
