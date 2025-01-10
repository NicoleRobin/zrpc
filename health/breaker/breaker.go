package breaker

import (
	"context"
	"strings"
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

var (
	registry = sync.Map{}

	breakerCount atomic.Uint32

	updateLock sync.RWMutex

	definedBreakers = breakerInfo
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
	if breakerCount.Load() >= maxBreakerCount {
		panic(ErrTooManyBreakers)
	}

	updateLock.RLock()
	defer updateLock.RUnlock()

	b := newBreaker()
	if opt, ok := getBreakerOption(name); ok {
		b.renew(opt)
	}
	// #TODO: 设置工厂方法

	breakerCount.Add(1)
}

func getBreakerOption(ctx context.Context, name string) (opt Option, bool) {
	if opt, ok := definedBreakers.breakers[name]; ok {
		return opt, true
	}
	for _, v := range definedBreakers.temps {
		if v.isTemplate && strings.HasPrefix(name, v.Name) {
			return v, true
		}
	}
	return opt, false
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

func Get(ctx context.Context, name string) *Breaker {
	return nil
}

func loadBreaker(ctx context.Context, name string) (*Breaker, bool) {
	v, ok := registry.Load(name)
	if !ok {
		return nil, false
	}
	return v.(*Breaker), true
}
