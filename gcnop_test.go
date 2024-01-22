package filecache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopGarbageCollector(t *testing.T) {
	prepareGCTestFiles(t)

	gc := NewNopGarbageCollector()

	gc.OnInstanceInit()
	gc.OnOperation()

	assert.FileExists(t, "./testdata/gc/test1.cache")
	assert.FileExists(t, "./testdata/gc/test1.cache--meta")
	assert.FileExists(t, "./testdata/gc/test2.cache")
	assert.FileExists(t, "./testdata/gc/test2.cache--meta")
}
