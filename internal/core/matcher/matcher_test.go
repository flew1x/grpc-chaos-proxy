package matcher

import (
	"testing"

	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
)

func TestCompile_ValidRegex(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: "^Do.*"}

	matcher, err := Compile(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if matcher == nil {
		t.Fatal("expected matcher, got nil")
	}
}

func TestCompile_InvalidRegex(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: "["}

	_, err := Compile(m)
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestMatcher_Match_ServiceAndMethod(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: "^Do.*"}
	matcher, _ := Compile(m)

	frame := &engine.Frame{Service: "TestSvc", Method: "DoSomething"}
	if !matcher.Match(frame) {
		t.Error("expected match for correct service and method")
	}
}

func TestMatcher_Match_ServiceCaseInsensitive(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: "^Do.*"}
	matcher, _ := Compile(m)

	frame := &engine.Frame{Service: "testsvc", Method: "DoSomething"}

	if !matcher.Match(frame) {
		t.Error("expected match for service case-insensitive")
	}
}

func TestMatcher_Match_MethodNoMatch(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: "^Do.*"}
	matcher, _ := Compile(m)

	frame := &engine.Frame{Service: "TestSvc", Method: "Other"}

	if matcher.Match(frame) {
		t.Error("expected no match for method")
	}
}

func TestMatcher_Match_ServiceNoMatch(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: "^Do.*"}
	matcher, _ := Compile(m)

	frame := &engine.Frame{Service: "OtherSvc", Method: "DoSomething"}

	if matcher.Match(frame) {
		t.Error("expected no match for service")
	}
}

func TestMatcher_Match_EmptyService(t *testing.T) {
	m := config.Match{Service: "", MethodRegex: "^Do.*"}
	matcher, _ := Compile(m)

	frame := &engine.Frame{Service: "AnySvc", Method: "DoSomething"}

	if !matcher.Match(frame) {
		t.Error("expected match when service is empty in matcher")
	}
}

func TestMatcher_Match_EmptyMethodRegex(t *testing.T) {
	m := config.Match{Service: "TestSvc", MethodRegex: ""}
	matcher, _ := Compile(m)

	frame := &engine.Frame{Service: "TestSvc", Method: "AnyMethod"}

	if !matcher.Match(frame) {
		t.Error("expected match when method regex is empty")
	}
}
