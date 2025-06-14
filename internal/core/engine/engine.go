package engine

import (
	"context"
	"errors"
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"

	"google.golang.org/grpc/metadata"
)

type Frame struct {
	Ctx     context.Context // context of the call, can be used to read metadata
	Service string          // service name, e.g., "com.example.Service"
	Method  string          // method name, e.g., "MethodName"
	MD      metadata.MD     // metadata of the call, can be used to read or modify metadata
}

// Injector introduces a fault – blocks, modifies, or interrupts a call
// The returned error will be propagated up and converted to a gRPC response
type Injector interface {
	Apply(f *Frame) error
}

var (
	registryMu sync.RWMutex
	registry   = map[entity.InjectorType]func(cfg any) (Injector, error){}
)

// Register is called in the init() of each adapter-injector
func Register(name entity.InjectorType, factory func(cfg any) (Injector, error)) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = factory
}

// BuildInjector creates an Injector based on the type and config
func BuildInjector(name entity.InjectorType, cfg any) (Injector, error) {
	return buildInjector(name, cfg)
}

func buildInjector(name entity.InjectorType, cfg any) (Injector, error) {
	registryMu.RLock()
	f := registry[name]
	registryMu.RUnlock()

	if f == nil {
		return nil, errors.New("unknown injector: " + name.String())
	}

	return f(cfg)
}

type compiledRule struct {
	serviceRE *regexp.Regexp
	methodRE  *regexp.Regexp
	inj       Injector
}

// Engine accepts a config loader and compiles rules.
type Engine struct {
	rulesMu sync.RWMutex
	rules   []*compiledRule
	cfgLdr  *config.Loader
	logger  *zap.Logger
}

// New creates an Engine and subscribes to hot-reload
func New(ldr *config.Loader, logger *zap.Logger) (*Engine, error) {
	e := &Engine{cfgLdr: ldr, logger: logger}

	if err := e.reload(); err != nil {
		return nil, err
	}

	// hot-reload goroutine
	go func() {
		for range ldr.Notify() {
			_ = e.reload() // ignore error – stay on old rules
		}
	}()

	return e, nil
}

// reload rebuilds runtime rules
func (e *Engine) reload() error {
	cfg := e.cfgLdr.Current()
	if cfg == nil {
		return apperr.ErrConfigNotLoaded
	}

	if len(cfg.Rules) == 0 {
		e.logger.Warn("[engine] no rules available, skipping reload")

		e.rulesMu.Lock()
		e.rules = nil
		e.rulesMu.Unlock()

		return nil
	}

	e.logger.Debug("[engine] reloading rules", zap.Int("count", len(cfg.Rules)))

	var compiled []*compiledRule

	for _, r := range cfg.Rules {
		if r.Disabled {
			continue
		}

		sRe := regexp.MustCompile("^" + regexp.QuoteMeta(strings.ToLower(r.Match.Service)) + "$")

		var (
			mRe *regexp.Regexp
			err error
		)

		if r.Match.MethodRegex != "" {
			mRe, err = regexp.Compile(r.Match.MethodRegex)
			if err != nil {
				e.logger.Error("[engine] failed to compile method_regex for rule %s: %v", zap.String("rule_name", r.Name), zap.Error(err))

				return err
			}
		}

		e.logger.Info("[engine] compiling rule", zap.String("rule_name", r.Name), zap.String("service", r.Match.Service), zap.String("method_regex", r.Match.MethodRegex))

		var inj Injector

		switch {
		case r.Action.Delay != nil:
			inj, err = buildInjector(entity.DelayType, r.Action.Delay)
			if err != nil {
				e.logger.Error("[engine] buildInjector error", zap.String("rule_name", r.Name), zap.Error(err))
			}
		case r.Action.Abort != nil:
			inj, err = buildInjector(entity.AbortType, r.Action.Abort)
			if err != nil {
				e.logger.Error("[engine] buildInjector error", zap.String("rule_name", r.Name), zap.Error(err))
			}
		case r.Action.Chaos != nil:
			inj, err = buildInjector(entity.ChaosType, r.Action.Chaos)
			if err != nil {
				e.logger.Error("[engine] buildInjector error", zap.String("rule_name", r.Name), zap.Error(err))
			}
		case r.Action.Spammer != nil:
			inj, err = buildInjector(entity.SpammerType, r.Action.Spammer)
			if err != nil {
				e.logger.Error("[engine] buildInjector error", zap.String("rule_name", r.Name), zap.Error(err))
			}
		}
		if inj == nil {
			e.logger.Info("[engine] skip rule: injector is nil", zap.String("rule_name", r.Name))

			continue
		}

		compiled = append(compiled, &compiledRule{serviceRE: sRe, methodRE: mRe, inj: inj})
	}

	e.rulesMu.Lock()
	e.rules = compiled
	e.rulesMu.Unlock()

	e.logger.Info("[engine] rules compiled", zap.Int("rules_count", len(compiled)))

	return nil
}

// Process iterates over rules, applies the first matching Injector
func (e *Engine) Process(f *Frame) error {
	e.rulesMu.RLock()
	rules := e.rules
	e.rulesMu.RUnlock()

	if len(rules) == 0 {
		e.logger.Warn("[engine] no rules available, returning error")

		return apperr.ErrNoMatchingRule
	}

	for _, c := range rules {
		e.logger.Debug("[engine] checking rule", zap.String("service", f.Service), zap.String("method", f.Method), zap.String("rule_service", c.serviceRE.String()), zap.String("rule_method_regex", c.methodRE.String()))

		if c.serviceRE.MatchString(strings.ToLower(f.Service)) && (c.methodRE == nil || c.methodRE.MatchString(f.Method)) {
			return c.inj.Apply(f)
		}
	}

	return apperr.ErrNoMatchingRule
}
