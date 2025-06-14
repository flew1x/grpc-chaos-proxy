package engine

import (
	"regexp"
	"testing"
)

type testInjector struct {
	called *bool
	retErr error
}

func (ti *testInjector) Apply(f *Frame) error {
	*ti.called = true

	return ti.retErr
}

func TestRegisterAndBuildInjector(t *testing.T) {
	called := false

	Register("test", func(cfg any) (Injector, error) {
		return &testInjector{called: &called}, nil
	})

	inj, err := buildInjector("test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inj == nil {
		t.Fatal("expected injector, got nil")
	}
}

func TestBuildInjectorUnknown(t *testing.T) {
	_, err := buildInjector("unknown", nil)
	if err == nil || err.Error() != "unknown injector: unknown" {
		t.Fatalf("expected unknown injector error, got: %v", err)
	}
}

func TestCompiledRuleMatch(t *testing.T) {
	called := false

	inj := &testInjector{called: &called}
	rule := &compiledRule{
		serviceRE: regexp.MustCompile(`^svc$`),
		methodRE:  regexp.MustCompile(`^mth$`),
		inj:       inj,
	}

	frame := &Frame{Service: "svc", Method: "mth"}
	if !rule.serviceRE.MatchString(frame.Service) || !rule.methodRE.MatchString(frame.Method) {
		t.Fatal("regexp match failed")
	}

	err := rule.inj.Apply(frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("injector was not called")
	}
}

func TestCompiledRuleNoMatch(t *testing.T) {
	inj := &testInjector{called: new(bool)}

	rule := &compiledRule{
		serviceRE: regexp.MustCompile(`^svc$`),
		methodRE:  regexp.MustCompile(`^mth$`),
		inj:       inj,
	}

	frame := &Frame{Service: "other", Method: "mth"}

	if rule.serviceRE.MatchString(frame.Service) {
		t.Fatal("should not match service")
	}
}
