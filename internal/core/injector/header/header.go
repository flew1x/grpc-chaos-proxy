package header

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
)

type ValueMod struct {
	Prefix string
	Suffix string
	Values []string
}

type Injector struct {
	Headers   map[string]ValueMod // key - header name, value - modification
	Direction entity.Direction    // "inbound", "outbound" or "both"
	Allowlist []string            // if not empty, only these headers will be kept
}

func NewHeaderInjector(cfg any) (engine.Injector, error) {
	conf, ok := cfg.(map[string]any)
	if !ok || conf == nil {
		return nil, apperr.ErrInvalidConfig
	}

	headers := make(map[string]ValueMod)

	if h, ok := conf["headers"].(map[string]any); ok {
		for k, v := range h {
			mod := ValueMod{}

			if m, ok := v.(map[string]any); ok {
				if prefix, ok := m["prefix"].(string); ok {
					mod.Prefix = prefix
				}

				if suffix, ok := m["suffix"].(string); ok {
					mod.Suffix = suffix
				}

				if vals, ok := m["values"].([]any); ok {
					for _, val := range vals {
						if s, ok := val.(string); ok {
							mod.Values = append(mod.Values, s)
						}
					}
				}
			} else if s, ok := v.(string); ok {
				mod.Values = []string{s}
			}

			headers[k] = mod
		}
	}

	direction := entity.DirectionBoth.String()

	if d, ok := conf["direction"].(string); ok {
		direction = d
	}

	var allowlist []string

	if a, ok := conf["allowlist"].([]any); ok {
		for _, v := range a {
			if s, ok := v.(string); ok {
				allowlist = append(allowlist, s)
			}
		}
	}

	return &Injector{
		Headers:   headers,
		Direction: entity.Direction(direction),
		Allowlist: allowlist,
	}, nil
}

func (h *Injector) Apply(f *engine.Frame) error {
	if f == nil || f.MD == nil {
		return nil
	}

	// direction check (inbound/outbound/both)
	if h.Direction != entity.DirectionBoth && string(f.Direction) != "" && f.Direction != h.Direction {
		return nil
	}

	// remove all headers except allowed
	if len(h.Allowlist) > 0 {
		for k := range f.MD {
			if !contains(h.Allowlist, k) {
				f.MD.Delete(k)
			}
		}
	}

	// add/modify headers
	for k, mod := range h.Headers {
		vals := make([]string, 0, len(mod.Values))

		if len(mod.Values) > 0 {
			for _, v := range mod.Values {
				vals = append(vals, mod.Prefix+v+mod.Suffix)
			}
		} else {
			vals = append(vals, mod.Prefix+mod.Suffix)
		}

		f.MD.Delete(k)

		for _, v := range vals {
			f.MD.Append(k, v)
		}
	}

	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}

func init() {
	engine.Register(entity.HeaderType, NewHeaderInjector)
}
