package filecache

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

// ScannerHitFn is a function called on every scanner's hit.
// The `key` parameter is a found cache item's key.
// The `createdAt` parameter is a time when item was created.
// The `opt` parameter is an item's ItemOptions data.
type ScannerHitFn func(key string, createdAt time.Time, opt *ItemOptions) error

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

		return onHit(meta.Key, meta.CreatedAt, metaToOptions(meta))
	})
}
