package filecache

import "time"

type GarbageCollector interface {
	OnInstanceInit()
	OnOperation()
	Close() error
}

// NewNopGarbageCollector returns the GarbageCollector doing nothing.
func NewNopGarbageCollector() GarbageCollector {
	return &gcNop{}
}

// NewProbabilityGarbageCollector returns the GarbageCollector running with the defined probability.
// Divisor is a run probability divisor (e.g., divisor equals 100 is a 1/100 probability).
func NewProbabilityGarbageCollector(dir string, onInitDivisor uint, onOpDivisor uint) GarbageCollector {
	return &gcProbability{
		dir:           dir,
		onInitDivisor: onInitDivisor,
		onOpDivisor:   onOpDivisor,
	}
}

// NewIntervalGarbageCollector returns the GarbageCollector running by the interval.
func NewIntervalGarbageCollector(dir string, interval time.Duration) GarbageCollector {
	return &gcInterval{
		dir:      dir,
		interval: interval,
	}
}
