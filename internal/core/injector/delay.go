package injector

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"math/rand"
	"time"
)

// DelayInjector sleeps for a random duration in [min,max] before forwarding the RPC
type DelayInjector struct {
	min time.Duration
	max time.Duration
}

// NewDelay builds the injector from config.DelayAction.
func NewDelay(cfg any) (engine.Injector, error) {
	dc, ok := cfg.(*config.DelayAction)
	if !ok || dc == nil {
		return nil, fmt.Errorf("delay action config error")
	}

	// fallback: if Max==0, use Min; if Min>Max, swap
	if dc.MaxMS == 0 {
		dc.MaxMS = dc.MinMS
	}

	if dc.MinMS > dc.MaxMS {
		dc.MinMS, dc.MaxMS = dc.MaxMS, dc.MinMS
	}

	return &DelayInjector{
		min: time.Duration(dc.MinMS) * time.Millisecond,
		max: time.Duration(dc.MaxMS) * time.Millisecond,
	}, nil
}

func (d *DelayInjector) Apply(f *engine.Frame) error {
	if d.max == 0 {
		return nil
	}

	dur := d.min

	if d.max > d.min {
		delta := d.max - d.min
		dur += time.Duration(rand.Int63n(int64(delta)))
	}

	select {
	case <-f.Ctx.Done():
		return f.Ctx.Err()
	case <-time.After(dur):
	}

	return nil
}

func init() {
	engine.Register(entity.DelayType, NewDelay)
}
