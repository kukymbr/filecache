// Copyright 2019 Sergey Basov. All right reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package filecache is a tool to write data from any io.Reader
// to cache files with TTL and metadata.
package filecache

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var (
	// NamespaceDefault is a default cache files namespace
	NamespaceDefault = "dft"

	// ExtDefault is a default files extension string
	ExtDefault = ".cache"

	// TTLDefault is a default value
	// in seconds of cache items' Time-To-Live
	TTLDefault = int64(-1)
)

// New creates new file cache instance
func New(path string) (*FileCache, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(path, 0744); err != nil {
		return nil, err
	}

	c := &FileCache{
		Path:             path,
		NamespaceDefault: NamespaceDefault,
		Ext:              ExtDefault,
		TTLDefault:       TTLDefault,
	}

	return c, nil
}

// FileCache is a file-based cache structure
type FileCache struct {
	// Path to the cache files directory
	Path string

	// NamespaceDefault is a default namespace if left empty in functions params
	NamespaceDefault string

	// Ext is a cache files extension string
	Ext string

	// TTLDefault is a default TTL of cache instance items
	TTLDefault int64
}

// Write copies data from src Reader to cache file
// Returns the count of written bytes
func (c *FileCache) Write(meta *Meta, src io.Reader) (written int64, err error) {
	item, written, err := c.WriteOpen(meta, src)
	if err != nil {
		return 0, err
	}
	_ = item.Close()
	return written, err
}

// WriteOpen copies data from src Reader to cache file
// and returns opened cache Item and count of written bytes
func (c *FileCache) WriteOpen(meta *Meta, src io.Reader) (item *Item, written int64, err error) {
	c.prepareMeta(meta)
	path, err := c.itemPath(meta.Key, meta.Namespace, false, true)
	if err != nil {
		return nil, 0, err
	}

	target, err := os.Create(path)
	if err != nil {
		return nil, 0, err
	}

	if err = c.writeMeta(path, meta); err != nil {
		_ = c.invalidatePath(path)
		return nil, 0, err
	}

	written, err = io.Copy(target, src)
	if err != nil {
		return nil, 0, err
	}

	item = &Item{
		File: target,
		Meta: meta,
		Path: path,
	}

	return item, written, nil
}

// Read returns cache Item if exists
func (c *FileCache) Read(key string, namespace string) (item *Item, err error) {
	path, err := c.itemPath(
		key,
		namespace,
		false,
		false,
	)
	if err != nil {
		return nil, err
	}

	meta := c.readMeta(path)
	if meta == nil {
		_ = c.invalidatePath(path)
		return nil, errors.New("failed to read meta for key" + key + " in namespace " + namespace)
	}
	if c.isExpired(meta) {
		_ = c.invalidatePath(path)
		return nil, errors.New("file is expired")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	item = &Item{
		File: f,
		Meta: meta,
		Path: path,
	}

	return item, nil
}

// Invalidate deletes cache item by its key & namespace
func (c *FileCache) Invalidate(key string, namespace string) error {
	path, err := c.itemPath(key, namespace, false, false)
	if err != nil {
		return err
	}
	return c.invalidatePath(path)
}

// invalidatePath deletes cache item by its path
func (c *FileCache) invalidatePath(itemPath string) error {
	var res error

	err := os.Remove(itemPath)
	if err != nil {
		res = err
	}

	path := c.metaFilePath(itemPath)
	err = os.Remove(path)
	if err != nil && res == nil {
		res = err
	}

	return res
}

// itemPath returns item's cache file path
func (c *FileCache) itemPath(key string, namespace string, relative bool, createDirs bool) (path string, err error) {
	key = c.itemKey(key)

	if namespace == "" {
		namespace = c.NamespaceDefault
	}

	dir := namespace + "/" + key[:2] + "/" + key[2:4] + "/" + key[4:6] + "/"
	dirAbs := c.Path + "/" + dir

	if !relative {
		dir = dirAbs
	}

	if createDirs {
		err := os.MkdirAll(dirAbs, 0744)
		if err != nil {
			return "", err
		}
	} else {
		if _, err := os.Stat(dirAbs); os.IsNotExist(err) {
			return "", errors.New("cache item directory does not exist: " + dirAbs)
		}
	}

	return dir + key[6:] + c.Ext, nil
}

// itemKey returns hex-encoded key hash string
func (c *FileCache) itemKey(key string) string {
	h := sha1.New()
	_, _ = io.WriteString(h, key)
	return hex.EncodeToString(h.Sum(nil))
}

// isExpired returns true if file is expired or if its TTL is 0
func (c *FileCache) isExpired(meta *Meta) bool {
	if meta.TTL == -1 {
		return false
	}
	now := time.Now().Unix()
	exp := meta.Created + meta.TTL
	return now > exp
}

// writeMeta data to file
func (c *FileCache) writeMeta(itemPath string, meta *Meta) error {
	meta.Created = time.Now().Unix()
	data, err := jsoniter.Marshal(meta)
	if err != nil {
		return err
	}
	path := c.metaFilePath(itemPath)
	return ioutil.WriteFile(path, data, 0744)
}

// readMeta data of cache item from file
// Returns nil if something goes wrong
func (c *FileCache) readMeta(itemPath string) *Meta {
	path := c.metaFilePath(itemPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	meta := &Meta{}
	if err = jsoniter.Unmarshal(data, &meta); err != nil {
		return nil
	}

	if meta.Key == "" {
		return nil
	}

	return meta
}

// metaFilePath returns the path of cache item metadata file
func (c *FileCache) metaFilePath(itemPath string) string {
	return itemPath + "--meta"
}

// prepareMeta sets default values to meta
func (c *FileCache) prepareMeta(meta *Meta) {
	if meta.Namespace == "" {
		meta.Namespace = c.NamespaceDefault
	}
	if meta.TTL == 0 {
		meta.TTL = c.TTLDefault
	}
}
