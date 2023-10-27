package filecache

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	// TTLEternal is a TTL value for eternal cache.
	TTLEternal = time.Duration(0)

	gcDivisorDefault = 100
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

	fc := &fileCache{
		dir:           targetDir,
		ttlDefault:    TTLEternal,
		gcDivisor:     gcDivisorDefault,
		pathGenerator: HashedKeySplitPath,
	}

	if len(options) == 1 {
		if options[0].DefaultTTL > 0 {
			fc.ttlDefault = options[0].DefaultTTL
		}

		if options[0].GCDivisor != 0 {
			fc.gcDivisor = options[0].GCDivisor
		}

		if options[0].PathGenerator != nil {
			fc.pathGenerator = options[0].PathGenerator
		}
	}

	return fc, nil
}

// FileCache is a tool to cache data from any io.Reader to the file.
type FileCache interface {
	// Write writes data from the reader to the cache file.
	Write(key string, reader io.Reader, options *ItemOptions) (written int, err error)

	// WriteData writes data to the cache file.
	WriteData(key string, data []byte, options *ItemOptions) (written int, err error)

	// Open opens the reader with cached data.
	Open(key string) (result *OpenResult, err error)

	// Read reads data from the cache file.
	Read(key string) (result *ReadResult, err error)
}

type fileCache struct {
	dir           string
	pathGenerator PathGeneratorFn
	ttlDefault    time.Duration
	gcDivisor     uint
}

func (fc *fileCache) Write(key string, reader io.Reader, options *ItemOptions) (written int, err error) {
	//TODO implement me
	panic("implement me")
}

func (fc *fileCache) WriteData(key string, data []byte, options *ItemOptions) (written int, err error) {
	//TODO implement me
	panic("implement me")
}

func (fc *fileCache) Open(key string) (result *OpenResult, err error) {
	//TODO implement me
	panic("implement me")
}

func (fc *fileCache) Read(key string) (result *ReadResult, err error) {
	//TODO implement me
	panic("implement me")
}

func (fc *fileCache) getItemPath(key string, options *ItemOptions, forMeta bool) string {
	path := fc.pathGenerator(key, options)

	if forMeta {
		path += metaPostfix
	}

	return path
}
