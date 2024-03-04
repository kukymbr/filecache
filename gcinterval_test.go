package filecache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntervalGarbageCollector(t *testing.T) {
	prepareGCTestFiles(t)

	gc := NewIntervalGarbageCollector("./testdata/gc", 50*time.Millisecond)

	gc.OnInstanceInit()
	gc.OnOperation()

	assert.FileExists(t, "./testdata/gc/test1.cache")
	assert.FileExists(t, "./testdata/gc/test1.cache--meta")
	assert.FileExists(t, "./testdata/gc/test2.cache")
	assert.FileExists(t, "./testdata/gc/test2.cache--meta")

	time.Sleep(55 * time.Millisecond)

	assert.NoFileExists(t, "./testdata/gc/test1.cache")
	assert.NoFileExists(t, "./testdata/gc/test1.cache--meta")
	assert.FileExists(t, "./testdata/gc/test2.cache")
	assert.FileExists(t, "./testdata/gc/test2.cache--meta")

	err := gc.Close()

	assert.NoError(t, err)
}
