package filecache_test

import (
	"gitlab.com/kukymbrgo/filecache"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileCache_WriteRead(t *testing.T) {
	cachePath, err := ioutil.TempDir("", "kukymbrgo-filecache-test")
	if err != nil {
		panic(err)
	}

	filecache.TTLDefault = 3600

	fc, err := filecache.New(cachePath)
	if err != nil {
		t.Error("failed to create filecache instance:", err)
		return
	}

	key := "testkey"
	sample := `test data string`
	sLen := len(sample)
	reader := strings.NewReader(sample)

	c, err := fc.Write(&filecache.Meta{Key: key}, reader)
	if err != nil {
		t.Error("failed to write data to cache:", err)
		return
	}

	if c != int64(sLen) {
		t.Error("written and expected string lengths does not match: expected", sLen, "got", c)
		return
	}

	item, err := fc.Read(key, "")
	if err != nil {
		t.Error("failed to read data from cache:", err)
		return
	}

	if item.Meta.Namespace != filecache.NamespaceDefault {
		t.Error("expected default namespace", filecache.NamespaceDefault, "got", item.Meta.Namespace)
	}

	if item.Meta.TTL != filecache.TTLDefault {
		t.Error("expected default TTL", filecache.TTLDefault, "got", item.Meta.TTL)
	}

	ext := filepath.Ext(item.Path)
	if ext != filecache.ExtDefault {
		t.Error("expected default extension", filecache.ExtDefault, "got", ext)
	}

	err = fc.Invalidate(key, "")
	if err != nil {
		t.Error("failed to invalidate cache item")
	}

	if err = os.RemoveAll(cachePath); err != nil {
		t.Error("failed to clean up after test")
	}
}
