package disconnect

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
)

type Disconnect struct {
	Percentage int `mapstructure:"percentage" yaml:"percentage"` // Percentage of requests to disconnect
}

// NewDisconnect builds the disconnect injector from config.DisconnectAction
func NewDisconnect(cfg any) (engine.Injector, error) {
	conf, ok := cfg.(*Disconnect)
	if ok != true || conf == nil {
		return nil, apperr.ErrInvalidConfig
	}

	if conf.Percentage < 0 || conf.Percentage > 100 {
		return nil, fmt.Errorf("percentage must be between 0 and 100")
	}

	return &Disconnect{
		Percentage: conf.Percentage,
	}, nil
}

func (d Disconnect) Apply(f *engine.Frame) error {
	if d.Percentage <= 0 {
		return nil
	}

	if rand.Intn(100) < d.Percentage {
		return status.Error(codes.Unavailable, "chaos disconnect injected")
	}

	return nil
}

func init() {
	engine.Register(entity.DisconnectType, NewDisconnect)
}
