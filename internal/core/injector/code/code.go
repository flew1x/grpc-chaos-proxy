package code

import (
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/injector/utils"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"strings"
	"time"
)

type CodeInjector struct {
	Code        string            `mapstructure:"code" yaml:"code"`                       // gRPC code to return
	Message     string            `mapstructure:"message" yaml:"message"`                 // Custom error message
	Percentage  int               `mapstructure:"percentage" yaml:"percentage"`           // Probability to inject error (0-100)
	Metadata    map[string]string `mapstructure:"metadata" yaml:"metadata"`               // Metadata to add to response
	DelayMS     int               `mapstructure:"delay_ms" yaml:"delay_ms"`               // Delay before returning error (ms)
	OnlyOn      []string          `mapstructure:"only_on_methods" yaml:"only_on_methods"` // Methods to apply to
	RepeatCount int               `mapstructure:"repeat_count" yaml:"repeat_count"`       // How many times to repeat error
}

// NewCodeInjector builds the code injector from config.CodeAction
func NewCodeInjector(cfg any) (engine.Injector, error) {
	conf, ok := cfg.(*CodeInjector)
	if !ok || conf == nil {
		return nil, apperr.ErrInvalidConfig
	}

	if conf.Percentage < 0 || conf.Percentage > 100 {
		return nil, fmt.Errorf("percentage must be between 0 and 100")
	}

	return conf, nil
}

func (c *CodeInjector) Apply(f *engine.Frame) error {
	if c.Percentage == 0 {
		return nil
	}

	if len(c.OnlyOn) > 0 {
		matched := false
		for _, m := range c.OnlyOn {
			if strings.EqualFold(m, f.Method) {
				matched = true
				break
			}
		}

		if !matched {
			return nil
		}
	}

	if c.Percentage > 0 && rand.Intn(100) >= c.Percentage {
		return nil
	}

	if c.RepeatCount > 0 {
		key := "x-chaos-repeat-" + c.Code
		count := 0

		if vals := f.MD.Get(key); len(vals) > 0 {
			fmt.Sscanf(vals[0], "%d", &count)
		}

		if count >= c.RepeatCount {
			return nil
		}

		count++
		f.MD.Set(key, fmt.Sprintf("%d", count))
	}

	if c.DelayMS > 0 {
		time.Sleep(time.Duration(c.DelayMS) * time.Millisecond)
	}

	if c.Metadata != nil && f.MD != nil {
		for k, v := range c.Metadata {
			f.MD.Set(k, v)
		}
	}

	grpcCode := codes.Unknown

	if code, ok := utils.CodeMap[strings.ToUpper(c.Code)]; ok {
		grpcCode = code
	}

	msg := c.Message
	if msg == "" {
		msg = "chaos code injected"
	}

	return status.Error(grpcCode, msg)
}

func init() {
	engine.Register(entity.CodeType, NewCodeInjector)
}
