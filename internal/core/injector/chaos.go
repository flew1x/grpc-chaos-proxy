package injector

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"math/rand"
	"time"

	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
)

// ChaosInjector randomly applies one of the configured actions
type ChaosInjector struct {
	injectors []engine.Injector
}

func NewChaosInjector(cfg any) (engine.Injector, error) {
	chaosCfg, ok := cfg.(*config.ChaosAction)
	if !ok || chaosCfg == nil || len(chaosCfg.Actions) == 0 {
		return nil, fmt.Errorf("invalid or empty config for ChaosInjector")
	}

	var injectors []engine.Injector

	for _, action := range chaosCfg.Actions {
		inj, err := buildInjectorFromAction(action)
		if err != nil {
			return nil, fmt.Errorf("failed to build injector: %w", err)
		}

		injectors = append(injectors, inj)
	}

	return &ChaosInjector{injectors: injectors}, nil
}

func (c *ChaosInjector) Apply(f *engine.Frame) error {
	if len(c.injectors) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(c.injectors))

	return c.injectors[idx].Apply(f)
}

// buildInjectorFromAction builds an injector based on the action type
func buildInjectorFromAction(action config.Action) (engine.Injector, error) {
	switch {
	case action.Delay != nil:
		return engine.BuildInjector(entity.DelayType, action.Delay)
	case action.Abort != nil:
		return engine.BuildInjector(entity.AbortType, action.Abort)
	case action.Spammer != nil:
		return engine.BuildInjector(entity.SpammerType, action.Spammer)
	default:
		return nil, fmt.Errorf("unknown action type in ChaosInjector")
	}
}

func init() {
	engine.Register(entity.ChaosType, NewChaosInjector)
}
