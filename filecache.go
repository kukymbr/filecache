package filecache

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kukymbr/filecache/v2/internal/util"
)

const (
	// TTLEternal is a TTL value for eternal cache.
	TTLEternal = util.TTLEternal
)

// New creates a new FileCache instance with a specified target dir & options.
func New(targetDir string, options ...InstanceOptions) (FileCache, error) {
	if len(options) > 1 {
		return nil, fmt.Errorf("more than one instance options param behavior is not supported")
	}

	targetDir = util.FixSeparators(targetDir)

	if targetDir == "" {
		targetDir = os.TempDir()
	}

	if err := util.PrepareDir(targetDir); err != nil {
		return nil, err
	}

	fc := &fileCache{
		dir:           targetDir,
		ttlDefault:    TTLEternal,
		pathGenerator: HashedKeySplitPath,
		gc:            NewProbabilityGarbageCollector(targetDir, 1, 100),
		keysLocker:    newKeysLocker(),
	}

	if len(options) == 1 {
		if options[0].DefaultTTL != 0 {
			fc.ttlDefault = options[0].DefaultTTL
		}

		if options[0].GC != nil {
			fc.gc = options[0].GC
		} else if options[0].GCDivisor != 0 {
			fc.gc = NewProbabilityGarbageCollector(targetDir, 1, options[0].GCDivisor)
		}

		if options[0].PathGenerator != nil {
			fc.pathGenerator = util.PathGeneratorFn(options[0].PathGenerator)
		}
	}

	go fc.gc.OnInstanceInit()

	return fc, nil
}

// NewInTemp creates a new FileCache instance with files stored in the system's temp dir.
func NewInTemp(options ...InstanceOptions) (FileCache, error) {
	return New("", options...)
}

// FileCache is a tool to cache data from any io.Reader to the file.
type FileCache interface {
	// GetPath returns the target path of the FileCache instance.
	GetPath() string

	// Write writes data from the reader to the cache file.
	Write(ctx context.Context, key string, reader io.Reader, options ...ItemOptions) (written int64, err error)

	// WriteData writes data to the cache file.
	WriteData(ctx context.Context, key string, data []byte, options ...ItemOptions) (written int64, err error)

	// Open opens the reader with cached data.
	//
	// Returns no error on successful cache hit, on no hit, on invalid cache files.
	// Returns an error if failed to open an existing cache file or if context is done.
	Open(ctx context.Context, key string) (result *OpenResult, err error)

	// Read reads data from the cache file.
	//
	// Returns no error on successful cache hit, on no hit, on invalid cache files.
	// Returns an error if failed to open or read an existing cache file or if context is done.
	Read(ctx context.Context, key string) (result *ReadResult, err error)

	// Invalidate removes data associated with a key from a cache.
	Invalidate(ctx context.Context, key string) error

	// Close closes the FileCache instance.
	Close() error
}

type fileCache struct {
	dir           string
	pathGenerator util.PathGeneratorFn
	ttlDefault    time.Duration
	gc            GarbageCollector

	keysLocker *keysLocker
}

func (fc *fileCache) GetPath() string {
	return fc.dir
}

//nolint:funlen
func (fc *fileCache) Write(
	ctx context.Context,
	key string,
	reader io.Reader,
	options ...ItemOptions,
) (written int64, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	defer func() {
		go fc.gc.OnOperation()
	}()

	opt := ItemOptions{}

	if len(options) > 0 {
		opt = options[0]
	}

	fc.keysLocker.lock(key)
	defer fc.keysLocker.unlock(key)

	meta := newMeta(key, &opt, fc.ttlDefault)
	itemPath := fc.getItemPath(key, false, true)
	metaPath := fc.getItemPath(key, true, true)

	itemF, err := create(key, itemPath)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = itemF.Close()
	}()

	metaF, err := create(key, metaPath)
	if err != nil {
		_ = itemF.Close()

		util.DeleteCacheFiles(itemPath, "")

		return 0, err
	}

	defer func() {
		_ = metaF.Close()
	}()

	undo := func() {
		_ = itemF.Close()
		_ = metaF.Close()

		util.DeleteCacheFiles(itemPath, metaPath)
	}

	if err := saveMeta(ctx, meta, metaF); err != nil {
		undo()

		return 0, err
	}

	n, err := copyWithCtx(ctx, itemF, reader)
	if err != nil {
		undo()

		return 0, err
	}

	return n, nil
}

func (fc *fileCache) WriteData(
	ctx context.Context,
	key string,
	data []byte,
	options ...ItemOptions,
) (written int64, err error) {
	reader := bytes.NewReader(data)

	return fc.Write(ctx, key, reader, options...)
}

func (fc *fileCache) Open(ctx context.Context, key string) (result *OpenResult, err error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	defer func() {
		go fc.gc.OnOperation()
	}()

	result = &OpenResult{}

	fc.keysLocker.lock(key)
	defer fc.keysLocker.unlock(key)

	itemPath := fc.getItemPath(key, false, false)
	metaPath := fc.getItemPath(key, true, false)

	if !util.ItemFilesValid(itemPath, metaPath) {
		util.DeleteCacheFiles(itemPath, metaPath)

		return result, nil
	}

	meta, err := readMeta(key, metaPath)
	if err != nil {
		util.DeleteCacheFiles(itemPath, metaPath)

		return result, nil
	}

	if meta.isExpired() {
		util.DeleteCacheFiles(itemPath, metaPath)

		return result, nil
	}

	result.hit = true
	result.options = metaToOptions(meta)

	result.reader, err = os.Open(itemPath)
	if err != nil {
		util.DeleteCacheFiles(itemPath, metaPath)

		return nil, fmt.Errorf("failed to open cache file for key %s: %w", key, err)
	}

	return result, nil
}

func (fc *fileCache) Read(ctx context.Context, key string) (result *ReadResult, err error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	defer func() {
		go fc.gc.OnOperation()
	}()

	itemPath := fc.getItemPath(key, false, false)
	metaPath := fc.getItemPath(key, true, false)

	openRes, err := fc.Open(ctx, key)
	if err != nil {
		return nil, err
	}

	result = &ReadResult{}

	if !openRes.Hit() {
		return result, nil
	}

	data, err := readAll(ctx, openRes.reader)
	if err != nil {
		util.DeleteCacheFiles(itemPath, metaPath)

		return nil, fmt.Errorf("failed to read cache data for key %s: %w", key, err)
	}

	result.hit = true
	result.options = openRes.options
	result.data = data

	return result, nil
}

func (fc *fileCache) Invalidate(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	defer func() {
		go fc.gc.OnOperation()
	}()

	fc.keysLocker.lock(key)
	defer fc.keysLocker.unlock(key)

	itemPath := fc.getItemPath(key, false, false)
	metaPath := fc.getItemPath(key, true, false)

	util.DeleteCacheFiles(itemPath, metaPath)

	return nil
}

func (fc *fileCache) Close() error {
	if err := fc.gc.Close(); err != nil {
		return err
	}

	return nil
}

func (fc *fileCache) getItemPath(key string, forMeta bool, createDirs bool) string {
	return util.GetItemPath(fc.GetPath(), fc.pathGenerator, key, forMeta, createDirs)
}
