package filecache

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// readItemMeta data of cache item from file
// Returns nil if something goes wrong
func readItemMeta(itemPath string) *Meta {
	path := itemToMetaPath(itemPath)
	meta, err := MetaFromFile(path)
	if err != nil {
		return nil
	}
	return meta
}

// itemToMetaPath returns the path of cache item metadata file
func itemToMetaPath(itemPath string) string {
	return itemPath + MetaPostfix
}

// pathIsMeta returns true is file name is meta file name
func pathIsMeta(path string) bool {
	lp := len(path)
	lm := len(MetaPostfix)
	if lp < lm {
		return false
	}
	return path[lp-lm:] == MetaPostfix
}

// itemKey returns hex-encoded key hash string
func itemKey(key string) string {
	h := sha1.New()
	_, _ = io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}

// invalidatePath deletes cache item by its path
func invalidatePath(itemPath string) error {
	var res error

	err := os.Remove(itemPath)
	if err != nil {
		res = err
	}

	path := itemToMetaPath(itemPath)
	err = os.Remove(path)
	if err != nil && res == nil {
		res = err
	}

	return res
}
