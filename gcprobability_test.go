package filecache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProbabilityGarbageCollector(t *testing.T) {
	prepareGCTestFiles(t)

	{
		gc := &gcProbability{}

		assert.False(t, gc.decideToRun(0))
	}

	{
		gc := &gcProbability{}

		assert.True(t, gc.decideToRun(1))
	}

	{
		gc := &gcProbability{
			dir: "./testdata/gc",
		}

		gc.run(1)

		assert.NoFileExists(t, "./testdata/gc/test1.cache")
		assert.NoFileExists(t, "./testdata/gc/test1.cache--meta")
		assert.FileExists(t, "./testdata/gc/test2.cache")
		assert.FileExists(t, "./testdata/gc/test2.cache--meta")
	}
}
