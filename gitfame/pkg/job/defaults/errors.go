package defaults

import "errors"

var (
	ErrAlreadyRegistered = errors.New("job is already registered")
	ErrJobNotRegistered  = errors.New("job not registered")
)
