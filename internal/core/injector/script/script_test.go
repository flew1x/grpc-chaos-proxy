package script

import (
	"context"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"strings"
	"testing"
)

func TestScriptInjector_Success(t *testing.T) {
	inj, err := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "exit 0",
		TimeoutMS: 500,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	frame := &engine.Frame{Ctx: context.Background()}

	if err := inj.Apply(frame); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestScriptInjector_Error(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "exit 1",
		TimeoutMS: 500,
	})

	frame := &engine.Frame{Ctx: context.Background()}

	err := inj.Apply(frame)
	if err == nil || !strings.Contains(err.Error(), "script error") {
		t.Errorf("expected script error, got %v", err)
	}
}

func TestScriptInjector_Timeout(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "sleep 2",
		TimeoutMS: 100,
	})

	frame := &engine.Frame{Ctx: context.Background()}

	err := inj.Apply(frame)
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestScriptInjector_Env(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "[ \"$FOO\" = \"bar\" ] || exit 2",
		Env:       map[string]string{"FOO": "bar"},
		TimeoutMS: 500,
	})

	frame := &engine.Frame{Ctx: context.Background()}

	if err := inj.Apply(frame); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestScriptInjector_Args(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "[ \"$1\" = \"baz\" ] || exit 3",
		Args:      []string{"baz"},
		TimeoutMS: 500,
	})

	frame := &engine.Frame{Ctx: context.Background()}

	if err := inj.Apply(frame); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestScriptInjector_UnsupportedLang(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language: "python",
		Source:   "print('hi')",
	})

	frame := &engine.Frame{Ctx: context.Background()}

	err := inj.Apply(frame)
	if err == nil || !strings.Contains(err.Error(), "only 'sh' or 'bash'") {
		t.Errorf("expected unsupported lang error, got %v", err)
	}
}

func TestScriptInjector_EmptySource(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language: "sh",
		Source:   "",
	})

	frame := &engine.Frame{Ctx: context.Background()}

	_ = inj.Apply(frame)
}

func TestScriptInjector_ChaosOutput_Error(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "echo 'X-CHAOS-ERROR: injected error'",
		TimeoutMS: 500,
	})

	frame := &engine.Frame{Ctx: context.Background()}

	err := inj.Apply(frame)
	if err == nil || !strings.Contains(err.Error(), "chaos script error: injected error") {
		t.Errorf("expected chaos script error, got %v", err)
	}
}

func TestScriptInjector_ChaosOutput_Header(t *testing.T) {
	inj, _ := NewScriptInjector(&config.ScriptAction{
		Language:  "sh",
		Source:    "echo 'X-CHAOS-HEADER: foo=bar'",
		TimeoutMS: 500,
	})

	frame := &engine.Frame{Ctx: context.Background(), MD: make(map[string][]string)}

	err := inj.Apply(frame)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if got := frame.MD["foo"]; len(got) == 0 || got[0] != "bar" {
		t.Errorf("expected header foo=bar, got %v", got)
	}
}
