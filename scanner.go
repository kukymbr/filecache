package filecache

import (
	"os"
	"path/filepath"
)

// NewScanner creates new Scanner instance
func NewScanner(fc *FileCache) *Scanner {
	return &Scanner{
		fc: fc,
	}
}

// Scanner is a tool to walk through existing cache files
type Scanner struct {
	fc *FileCache
}

// ScannerHitFunc is a function called on every cache file hit while scanning.
// Receives found cache item meta, path of cached content file & its info.
type ScannerHitFunc = func(meta *Meta, path string, info os.FileInfo) error

// Scan walks through existing cache files
// and executes the hit function on every cache file found.
func (s *Scanner) Scan(hitFunc ScannerHitFunc, skipExpired bool, ignoreLStatErrors bool) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if ignoreLStatErrors {
				return nil
			}
			// if file disappeared while scrolling
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info.IsDir() {
			return nil
		}
		if fileIsMeta(path) {
			return nil
		}
		if meta := readItemMeta(path); meta != nil {
			if skipExpired && meta.IsExpired() {
				return nil
			}
			return hitFunc(meta, path, info)
		}
		return nil
	}

	return filepath.Walk(s.fc.Path(), walkFn)
}
