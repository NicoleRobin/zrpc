package breaker

type EventOptions struct {
	Name      string
	Threshold int
	Interval  int
	Unit      int
}

type Option struct {
	Name      string
	Successes int
	Events    []EventOptions
	Enable    bool

	isTemplate bool
}

type breakerInfo struct {
	breakers map[string]Option
	temps    []Option
}

type breakerConfig struct {
	Breakers []Option
}
