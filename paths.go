package filecache

import (
	//nolint:gosec
	"crypto/sha1"
	"encoding/hex"
	"io"
	"path/filepath"
	"strings"

	"github.com/kukymbr/filecache/v2/internal/util"
)

// PathGeneratorFn is a function to generate cache item's file path.
// Receives the key of the cache item and returns the path of the item relative to the cache instance's dir.
type PathGeneratorFn util.PathGeneratorFn

// WithExt returns new PathGeneratorWithExt instance.
func WithExt(fn PathGeneratorFn, ext string) PathGeneratorFn {
	ext = util.FilterPathIdent(ext)

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return func(key string) string {
		return fn(key) + ext
	}
}

// FilteredKeyPath uses a key without path separators as a path.
func FilteredKeyPath(key string) string {
	path := util.FilterPathIdent(key)

	if path == "" {
		return HashedKeyPath(key)
	}

	return path
}

// HashedKeyPath return hashed key and uses it as a file name.
func HashedKeyPath(key string) string {
	//nolint:gosec
	h := sha1.New()
	_, _ = io.WriteString(h, key)

	return hex.EncodeToString(h.Sum(nil))
}

// HashedKeySplitPath return hashes key, splits it on the parts which are directories and a file name.
func HashedKeySplitPath(key string) string {
	//nolint:gosec
	h := sha1.New()
	_, _ = io.WriteString(h, key)

	hashed := hex.EncodeToString(h.Sum(nil))

	return filepath.Join(hashed[:2], hashed[2:4], hashed[4:6], hashed[6:])
}
