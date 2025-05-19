package filecache

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kukymbr/filecache/v2/internal/util"
	"github.com/mailru/easyjson"
)

//go:generate easyjson -omit_empty -no_std_marshalers -all meta.go

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

func (m meta) isExpired() bool {
	return util.IsExpired(m.CreatedAt, m.TTL)
}

func saveMeta(ctx context.Context, meta *meta, target *os.File) error {
	data, err := easyjson.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta for key %s: %w", meta.Key, err)
	}

	if _, err := copyWithCtx(ctx, target, bytes.NewReader(data)); err != nil {
		return fmt.Errorf("failed to save meta for key %s: %w", meta.Key, err)
	}

	return nil
}

func readMeta(key string, path string) (*meta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read meta file for key %s: %w", key, err)
	}

	var m meta

	if err := easyjson.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal meta for key %s: %w", key, err)
	}

	return &m, nil
}

func newMeta(key string, options *ItemOptions, defaultTTL time.Duration) *meta {
	ttl := defaultTTL

	if options.TTL != 0 {
		ttl = options.TTL
	}

	return &meta{
		Key:       key,
		CreatedAt: time.Now(),
		Name:      options.Name,
		TTL:       ttl,
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
