package breaker

import "sync"

type breakerState int

const (
	stateOpen breakerState = iota
	stateClosed
	stateHalfOpen

	maxBreakerCount = 1000
)

var breakerStateString = map[breakerState]string{
	stateOpen:     "open",
	stateClosed:   "closed",
	stateHalfOpen: "half-open",
}

func (s breakerState) String() string {
	return breakerStateString[s]
}

type Breaker struct {
	m sync.Mutex

	state breakerState

	enabled uint32
}
