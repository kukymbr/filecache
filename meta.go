package filecache

import (
	"time"
)

const metaPostfix = "--meta"

func newMeta(key string, options *ItemOptions) *meta {
	return &meta{
		Key:       key,
		CreatedAt: time.Now(),
		Name:      options.Name,
		TTL:       options.TTL,
		Fields:    options.Fields,
	}
}

func metaToOptions(meta *meta) *ItemOptions {
	return &ItemOptions{
		Name:   meta.Name,
		TTL:    meta.TTL,
		Fields: meta.Fields,
	}
}

// meta is a metadata stored with a cache item file.
type meta struct {
	// Key is a unique cache item key.
	Key string `json:"k"`

	// CreatedAt is a time when cache item was created.
	CreatedAt time.Time `json:"c"`

	// Name is a human-readable item name.
	Name string `json:"n,omitempty"`

	// TTL is an item's time-to-live value.
	TTL time.Duration `json:"t,omitempty"`

	// Fields is a map of any other metadata fields.
	Fields Values `json:"f,omitempty"`
}

func (m *meta) isExpired() bool {
	if m.TTL == TTLEternal || m.TTL < 0 {
		return false
	}

	return time.Since(m.CreatedAt) > m.TTL
}
