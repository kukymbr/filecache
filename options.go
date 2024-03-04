package filecache

import "time"

// InstanceOptions are a cache instance options.
type InstanceOptions struct {
	// PathGenerator is a function to generate cache item's file path.
	// Receives the key of the cache item and returns the path of the item relative to the cache instance's dir.
	//
	// There are some built-in path generators:
	// 	 - FilteredKeyPath: removes path separators from the key and uses it as a file name;
	//   - HashedKeyPath: hashes the key and uses the result as a file name;
	//   - HashedKeySplitPath: hashes the key, splits result to the parts and uses them as a directories and file names.
	//
	// Also, there is a WithExt wrapper adding the extension to any path generator result.
	PathGenerator PathGeneratorFn

	// DefaultTTL is a TTL value for the items without it.
	DefaultTTL time.Duration

	// GC is a GarbageCollector instance for the cache instance.
	//
	// May be initialized with any GarbageCollector instance or using one of the predefined GC constructors:
	//   - NewNopGarbageCollector: the GarbageCollector doing nothing;
	//   - NewProbabilityGarbageCollector: the GarbageCollector running with the defined probability;
	//   - NewIntervalGarbageCollector: the GarbageCollector running by the interval.
	GC GarbageCollector

	// GCDivisor is a garbage collector run probability divisor
	// (e.g., 100 is 1/100 probability).
	//
	// Deprecated: use the GC property instead.
	GCDivisor uint
}

// ItemOptions are a cache item options.
type ItemOptions struct {
	// Name is a human-readable item name.
	Name string

	// TTL is an item's time-to-live value.
	TTL time.Duration

	// Fields is a map of any other metadata fields.
	Fields Values
}
