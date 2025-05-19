package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PathGeneratorFn func(key string) string

// DeleteCacheFiles removes cache files
func DeleteCacheFiles(paths ...string) {
	if len(paths) > 2 {
		panic("unexpected behaviour: DeleteCacheFiles expects no more than two paths")
	}

	for _, path := range paths {
		_ = os.Remove(path)
	}
}

// PrepareDir checks if dir exists and creates it otherwise.
func PrepareDir(dir string) error {
	err := validateDir(dir)

	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrDirNotExists) {
		return err
	}

	if err = os.MkdirAll(dir, DirsMode); err != nil {
		return fmt.Errorf("%s dir does not exist and cannot be created: %w", dir, err)
	}

	return nil
}

// validateDir checks if a given path is an existing dir path.
func validateDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: %w", dir, ErrDirNotExists)
		}

		return fmt.Errorf("stat %s: %w", dir, err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s: %w", dir, ErrNotADir)
	}

	return nil
}

// ItemFilesValid checks if itemPath & metaPath are a valid files' paths.
func ItemFilesValid(itemPath string, metaPath string) bool {
	if itemPath == "" || metaPath == "" {
		return false
	}

	itemStat, err := os.Stat(itemPath)
	if err != nil {
		return false
	}

	metaStat, err := os.Stat(metaPath)
	if err != nil {
		return false
	}

	return !itemStat.IsDir() && !metaStat.IsDir()
}

// FixSeparators replaces all path separators with the OS-correct.
func FixSeparators(path string) string {
	sepToReplace := '/'
	if os.PathSeparator == sepToReplace {
		sepToReplace = '\\'
	}

	return strings.ReplaceAll(path, string(sepToReplace), string(os.PathSeparator))
}

// FilterPathIdent remove path separators from the path part.
func FilterPathIdent(ident string) string {
	ident = strings.TrimSpace(ident)
	ident = strings.ReplaceAll(ident, "/", "")
	ident = strings.ReplaceAll(ident, "\\", "")

	return ident
}

// IsExpired checks if item is expired.
func IsExpired(createdAt time.Time, ttl time.Duration) bool {
	if ttl == TTLEternal || ttl <= 0 {
		return false
	}

	return time.Since(createdAt) > ttl
}

// GetItemPath returns full item's path.
func GetItemPath(dir string, pathGenerator PathGeneratorFn, key string, forMeta bool, createDirs bool) string {
	path := filepath.Join(dir, pathGenerator(key))
	itemDir := filepath.Dir(FixSeparators(path))

	if itemDir != "." && createDirs {
		_ = os.MkdirAll(itemDir, DirsMode)
	}

	if forMeta {
		path += MetaSuffix
	}

	return path
}
