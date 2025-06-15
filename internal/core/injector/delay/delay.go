package delay

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"math/rand"
	"time"
)

// Injector DelayInjector sleeps for a random duration in [min,max] before forwarding the RPC
type Injector struct {
	min time.Duration
	max time.Duration
}

// NewDelayInjector builds the injector from config.DelayAction.
func NewDelayInjector(cfg any) (engine.Injector, error) {
	dc, ok := cfg.(*config.DelayAction)
	if !ok || dc == nil {
		return nil, apperr.ErrInvalidConfig
	}

	// fallback: if Max==0, use Min; if Min>Max, swap
	if dc.MaxMS == 0 {
		dc.MaxMS = dc.MinMS
	}

	if dc.MinMS > dc.MaxMS {
		dc.MinMS, dc.MaxMS = dc.MaxMS, dc.MinMS
	}

	return &Injector{
		min: time.Duration(dc.MinMS) * time.Millisecond,
		max: time.Duration(dc.MaxMS) * time.Millisecond,
	}, nil
}

func (d *Injector) Apply(f *engine.Frame) error {
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
	engine.Register(entity.DelayType, NewDelayInjector)
}
