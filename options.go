package filecache

import "time"

// InstanceOptions are a cache instance options.
type InstanceOptions struct {
	// FilesExt is a generated cache files extension.
	FilesExt string

	// NoExtensions is a flag to disable files extensions at all.
	NoExtensions bool

	// DefaultTTL is a TTL value for the items without it.
	DefaultTTL time.Duration

	// GCDivisor is a garbage collector run probability divisor
	// (e.g., 100 is 1/100 probability).
	GCDivisor uint
}

// ItemOptions are a cache item options
type ItemOptions struct {
	// Name is a human-readable item name.
	Name string

	// TTL is an item's time-to-live value.
	TTL time.Duration

	// Fields is a map of any other metadata fields.
	Fields Values
}
