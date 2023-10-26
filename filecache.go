package filecache

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	// TTLEternal is a TTL value for eternal cache.
	TTLEternal = time.Duration(-1)

	// FilesExtDefault is a default generated files extension.
	FilesExtDefault = ".cache"
)

// New creates new file cache instance with a specified target dir & options.
func New(targetDir string, options ...InstanceOptions) (FileCache, error) {
	if len(options) > 1 {
		return nil, fmt.Errorf("more than one instance options param behavior is not supported")
	}

	targetDir = fixSeparators(targetDir)

	if targetDir == "" {
		targetDir = os.TempDir()
	}

	if err := prepareDir(targetDir); err != nil {
		return nil, err
	}

	fc := &fileCache{dir: targetDir}

	if len(options) == 1 {
		if options[0].FilesExt != "" {
			fc.filesExt = options[0].FilesExt
		} else if !options[0].NoExtensions {
			fc.filesExt = FilesExtDefault
		}

		if options[0].DefaultTTL != 0 {
			fc.ttlDefault = options[0].DefaultTTL
		} else {
			fc.ttlDefault = TTLEternal
		}

		if options[0].GCDivisor != 0 {
			fc.gcDivisor = options[0].GCDivisor
		}
	}

	return fc, nil
}

type FileCache interface {
	Write(key string, reader io.Reader) (written int, err error)
	Read(key string) (reader io.ReadCloser, err error)
}

type fileCache struct {
	dir        string
	filesExt   string
	ttlDefault time.Duration
	gcDivisor  uint
}

func (f *fileCache) Write(key string, reader io.Reader) (written int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *fileCache) Read(key string) (reader io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}
