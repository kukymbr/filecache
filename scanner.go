package filecache

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// ScanEntry is a scanner hit entry.
type ScanEntry struct {
	// Key is a cache item key.
	Key string

	// CreatedAt is a cache item created-at timestamp.
	CreatedAt time.Time

	// Options are the options of the item stored in the cache.
	Options *ItemOptions

	itemPath string
	metaPath string
}

// ScannerHitFn is a function called on every scanner's hit.
// Function receives the ScanEntry, describing the found cache item.
// If the function returns an error, the iteration will be stopped.
type ScannerHitFn func(entry ScanEntry) error

// NewScanner creates a Scanner looking for the valid cache items.
func NewScanner(dir string) Scanner {
	return &scanner{dir: dir}
}

// newExpiredScanner creates a Scanner looking for the expired items.
func newExpiredScanner(dir string) Scanner {
	return &scanner{
		dir:         dir,
		expiredOnly: true,
	}
}

// Scanner is a tool to scan cache items inside the specified directory.
type Scanner interface {
	Scan(onHit ScannerHitFn) error
}

type scanner struct {
	dir         string
	expiredOnly bool
}

func (s *scanner) Scan(onHit ScannerHitFn) error {
	return filepath.WalkDir(s.dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if !strings.HasSuffix(entry.Name(), metaSuffix) {
			return nil
		}

		itemPath := strings.TrimSuffix(path, metaSuffix)
		metaPath := path

		if !itemFilesValid(itemPath, metaPath) {
			return nil
		}

		meta, err := readMeta("", path)
		if err != nil {
			//nolint:nilerr
			return nil
		}

		if s.expiredOnly && !meta.isExpired() || !s.expiredOnly && meta.isExpired() {
			return nil
		}

		return onHit(ScanEntry{
			Key:       meta.Key,
			CreatedAt: meta.CreatedAt,
			Options:   metaToOptions(meta),
			itemPath:  itemPath,
			metaPath:  metaPath,
		})
	})
}
