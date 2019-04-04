package filecache_test

import (
	"gitlab.com/kukymbrgo/filecache"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGC(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping GC test")
	}

	cachePath, err := ioutil.TempDir("", "kukymbrgo-filecache-test")
	if err != nil {
		panic(err)
	}

	filecache.TTLDefault = 1
	filecache.GCDivisor = 0

	fc, err := filecache.New(cachePath)
	if err != nil {
		t.Error("failed to create filecache instance:", err)
		return
	}

	defer func() {
		if err = os.RemoveAll(cachePath); err != nil {
			t.Error("failed to clean up after test")
		}
	}()

	for i := 0; i < 10; i++ {
		key := strconv.FormatInt(int64(i), 10)
		reader := strings.NewReader("test string #" + key)
		_, err = fc.Write(&filecache.Meta{Key: key}, reader)
		if err != nil {
			t.Error("failed to write #" + key)
		}
	}

	key := "10"
	reader := strings.NewReader("test string #" + key)
	_, err = fc.Write(&filecache.Meta{Key: key, TTL: 100}, reader)
	if err != nil {
		t.Error("failed to write #" + key)
	}

	time.Sleep(2 * time.Second)

	filecache.GCDivisor = 1

	fc, err = filecache.New(cachePath)
	if err != nil {
		t.Error("failed to recreate filecache instance:", err)
		return
	}

	count, err := countFiles(fc.Path())
	if err != nil {
		t.Error("failed to count files in dir", fc.Path())
		return
	}

	if count != 2 {
		t.Error("expected 2 files after GC run, got", count)
	}
}

// countFiles in dir recursively
func countFiles(path string) (count int, err error) {
	count = 0
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		count++
		return nil
	}
	err = filepath.Walk(path, walkFn)
	return
}
