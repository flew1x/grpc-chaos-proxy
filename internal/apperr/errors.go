package apperr

import "errors"

var (
	ErrConfigNotLoaded   = errors.New("config is not loaded")
	ErrNoMatchingRule    = errors.New("no matching rule found")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
