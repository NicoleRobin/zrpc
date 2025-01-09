package breaker

import (
	"sync"
	"sync/atomic"
)

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

	enabled atomic.Bool
}

func newBreaker() *Breaker {
	b := &Breaker{
		state: stateClosed,
	}

	return b
}

func initBreaker(name string) {

}

func (b *Breaker) Enable() {
	b.enabled.Store(true)
}

func (b *Breaker) Disable() {
	b.enabled.Store(false)
}

func (b *Breaker) Do(f func() error) error {
	return nil
}
