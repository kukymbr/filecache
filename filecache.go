package filecache

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// Default namespace
const NamespaceDefault = "dft"

// Create file cache instance
func New(path string, namespaceDefault string) (*FileCache, error) {
	err := os.MkdirAll(path, 0744)
	if err != nil {
		return nil, err
	}

	if namespaceDefault == "" {
		namespaceDefault = NamespaceDefault
	}

	c := &FileCache{
		Path:             path,
		NamespaceDefault: namespaceDefault,
	}

	return c, nil
}

// Cache
type FileCache struct {

	// Path to the cache files directory
	Path string

	// Default namespace if left empty in functions params
	NamespaceDefault string
}

// Write data from reader to the cache item
func (c *FileCache) WriteFromReader(key string, source io.Reader, namespace string) error {
	path, err := c.ItemPath(key, namespace)
	if err != nil {
		return err
	}

	target, err := os.Create(path)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := target.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}

// Get cache item reader
func (c *FileCache) Reader(key string, namespace string) (io.Reader, error) {
	path, err := c.ItemPath(key, namespace)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Get Key path
func (c *FileCache) ItemPath(key string, namespace string) (string, error) {
	if namespace == "" {
		namespace = c.NamespaceDefault
	}
	key = c.Key(key)
	dir := c.Path
	dir += "/" + namespace + "/" + key[:2] + "/" + key[2:4] + "/" + key[4:6] + "/"

	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return "", err
	}

	return dir + key[6:] + ".cache", nil
}

// Get hex-encoded Key hash string
func (c *FileCache) Key(key string) string {
	h := sha1.New()
	io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}
