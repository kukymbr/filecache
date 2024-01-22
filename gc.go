package filecache

import "time"

// GarbageCollector is a tool to remove expired cache items.
type GarbageCollector interface {
	// OnInstanceInit is executed on the initialization of the FileCache instance.
	OnInstanceInit()

	// OnOperation is executed on the every item's operation in the FileCache instance.
	OnOperation()

	// Close closes the GarbageCollector.
	Close() error
}

// NewNopGarbageCollector returns the GarbageCollector doing nothing.
func NewNopGarbageCollector() GarbageCollector {
	return &gcNop{}
}

// NewProbabilityGarbageCollector returns the GarbageCollector running with the defined probability.
// Divisor is a run probability divisor (e.g., divisor equals 100 is a 1/100 probability).
//
// Function arguments:
// * dir - the directory with the FileCache's instance files;
// * onInitDivisor - divisor for the probability on the OnInstanceInit() function call;
// * onOpDivisor - divisor for the probability on the OnOperation() function call.
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
