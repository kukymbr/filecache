package filecache

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const dirsMode os.FileMode = 0755

var (
	ErrDirNotExists = errors.New("directory does not exist")
	ErrNotADir      = errors.New("not a directory")
)

// prepareDir checks if dir exists and creates it otherwise.
func prepareDir(dir string) error {
	err := validateDir(dir)

	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrDirNotExists) {
		return err
	}

	if err = os.MkdirAll(dir, dirsMode); err != nil {
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

// itemFilesValid checks if itemPath & metaPath are a valid files' paths.
func itemFilesValid(itemPath string, metaPath string) bool {
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

// fixSeparators replaces all path separators with the OS-correct.
func fixSeparators(path string) string {
	sepToReplace := '/'
	if os.PathSeparator == sepToReplace {
		sepToReplace = '\\'
	}

	return strings.ReplaceAll(path, string(sepToReplace), string(os.PathSeparator))
}

func invalidate(itemPath string, metaPath string) {
	if itemPath != "" {
		_ = os.Remove(itemPath)
	}

	if metaPath != "" {
		_ = os.Remove(metaPath)
	}
}
