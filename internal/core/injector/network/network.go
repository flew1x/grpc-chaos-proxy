package network

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/injector/utils"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"time"
)

// Injector NetworkInjector is an injector that simulates network issues such as packet loss,
// reordering, and throttling. It can be configured with various parameters
type Injector struct {
	LossPercentage int `mapstructure:"loss_percentage" yaml:"loss_percentage"`
	ThrottleMS     int `mapstructure:"throttle_ms" yaml:"throttle_ms"`
}

// NewNetworkInjector creates a new NetworkInjector instance from the provided configuration
func NewNetworkInjector(cfg any) (engine.Injector, error) {
	nc, ok := cfg.(*Injector)
	if !ok {
		return nil, apperr.ErrInvalidConfig
	}

	if nc.LossPercentage < 0 || nc.LossPercentage > 100 {
		return nil, fmt.Errorf("loss percentage must be between 0 and 100")
	}

	if nc.ThrottleMS < 0 {
		return nil, fmt.Errorf("throttle ms must be non-negative")
	}

	return nc, nil
}

// Apply implements the engine.Injector interface for NetworkInjector
func (d *Injector) Apply(f *engine.Frame) error {
	if d.LossPercentage > 0 && utils.RandInt(100) < d.LossPercentage {
		return fmt.Errorf("network chaos: simulated packet loss (%d%%)", d.LossPercentage)
	}

	if d.ThrottleMS > 0 {
		time.Sleep(time.Duration(d.ThrottleMS) * time.Millisecond)
	}

	return nil
}

func init() {
	engine.Register(entity.NetworkType, NewNetworkInjector)
}
