package apperr

import "errors"

var (
	ErrConfigNotLoaded = errors.New("config is not loaded")
	ErrNoMatchingRule  = errors.New("no matching rule found")
)
