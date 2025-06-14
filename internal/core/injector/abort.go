package injector

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"math/rand"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	fullPercentage = 100 // represents 100% chance of aborting
)

// AbortInjector forcibly terminates an RPC with the configured gRPC status code
// in percentage % of cases
type AbortInjector struct {
	code       codes.Code
	percentage int // 0-100
}

// NewAbort builds the injector from config.AbortAction
func NewAbort(cfg any) (engine.Injector, error) {
	ac, ok := cfg.(*config.AbortAction)
	if !ok {
		return nil, fmt.Errorf("abort action config is not config.AbortAction")
	}

	code := codes.Internal

	if c, ok := codeMap[ac.Code]; ok {
		code = c
	}

	pct := ac.Percentage

	if pct < 0 {
		pct = 0
	}

	if pct > 100 {
		pct = 100
	}

	return &AbortInjector{code: code, percentage: pct}, nil
}

// Apply injects an error with the configured gRPC status code
func (ai *AbortInjector) Apply(*engine.Frame) error {
	if ai.percentage == 0 {
		return nil
	}

	if rand.Intn(fullPercentage) < ai.percentage {
		return status.Error(ai.code, "chaos abort injected")
	}

	return nil
}

func init() {
	engine.Register(entity.AbortType, NewAbort)
}
