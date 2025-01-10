package breaker

import (
	"errors"
)

var (
	ErrTooManyBreakers = errors.New("too many breakers")
)
