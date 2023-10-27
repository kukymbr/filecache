package filecache

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"path/filepath"
)

// PathGeneratorFn is a function to generate cache item's file path.
type PathGeneratorFn func(key string, options *ItemOptions) string

// HashedKeyPath return hashes key and uses it as a file name.
func HashedKeyPath(key string, _ *ItemOptions) string {
	h := sha1.New()
	_, _ = io.WriteString(h, key)

	return hex.EncodeToString(h.Sum(nil))
}

// HashedKeySplitPath return hashes key, splits it on the parts which are directories and a file name.
func HashedKeySplitPath(key string, _ *ItemOptions) string {
	h := sha1.New()
	_, _ = io.WriteString(h, key)

	hashed := hex.EncodeToString(h.Sum(nil))

	return filepath.Join(hashed[:2], hashed[2:4], hashed[4:6], hashed[6:])
}
