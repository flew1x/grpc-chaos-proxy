package header

import (
	"context"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"google.golang.org/grpc/metadata"
	"reflect"
	"testing"
)

func TestHeaderInjector_Apply_AddModifyHeaders(t *testing.T) {
	h := &Injector{
		Headers: map[string]ValueMod{
			"x-test": {Prefix: "pre-", Suffix: "-suf", Values: []string{"val1", "val2"}},
		},
		Direction: "both",
	}

	md := metadata.New(nil)
	frame := &engine.Frame{Ctx: context.Background(), MD: md, Direction: "both"}

	h.Apply(frame)
	got := frame.MD.Get("x-test")
	want := []string{"pre-val1-suf", "pre-val2-suf"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestHeaderInjector_Apply_Allowlist(t *testing.T) {
	h := &Injector{
		Headers:   map[string]ValueMod{},
		Allowlist: []string{"x-keep"},
		Direction: "both",
	}

	md := metadata.New(map[string]string{"x-keep": "1", "x-drop": "2"})
	frame := &engine.Frame{Ctx: context.Background(), MD: md, Direction: "both"}

	h.Apply(frame)

	if frame.MD.Get("x-keep")[0] != "1" {
		t.Error("x-keep should be kept")
	}

	if len(frame.MD.Get("x-drop")) != 0 {
		t.Error("x-drop should be deleted")
	}
}

func TestHeaderInjector_Apply_Direction(t *testing.T) {
	h := &Injector{
		Headers:   map[string]ValueMod{"x-dir": {Values: []string{"ok"}}},
		Direction: "inbound",
	}

	md := metadata.New(nil)
	frame := &engine.Frame{Ctx: context.Background(), MD: md, Direction: "outbound"}

	h.Apply(frame)

	if len(frame.MD.Get("x-dir")) != 0 {
		t.Error("header should not be injected for wrong direction")
	}

	frame.Direction = "inbound"

	h.Apply(frame)

	if frame.MD.Get("x-dir")[0] != "ok" {
		t.Error("header should be injected for correct direction")
	}
}

func TestHeaderInjector_Apply_Empty(t *testing.T) {
	h := &Injector{}

	frame := &engine.Frame{Ctx: context.Background(), MD: metadata.New(nil), Direction: "both"}

	if err := h.Apply(frame); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
