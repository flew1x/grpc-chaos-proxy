package matcher

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"regexp"
	"strings"
)

type Matcher struct {
	svcEq  string
	methRe *regexp.Regexp
}

func Compile(m config.Match) (*Matcher, error) {
	var (
		re  *regexp.Regexp
		err error
	)

	if m.MethodRegex != "" {
		re, err = regexp.Compile(m.MethodRegex)
		if err != nil {
			return nil, err
		}
	}

	return &Matcher{
		svcEq:  strings.ToLower(m.Service),
		methRe: re,
	}, nil
}

func (mt *Matcher) Match(f *engine.Frame) bool {
	if mt.svcEq != "" && strings.ToLower(f.Service) != mt.svcEq {
		return false
	}

	if mt.methRe != nil && !mt.methRe.MatchString(f.Method) {
		return false
	}

	return true
}
