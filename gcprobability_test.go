package filecache

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGarbageCollector(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("./testdata/gc")
	})

	{
		err := os.MkdirAll("./testdata/gc", dirsMode)
		require.NoError(t, err)

		err = os.WriteFile("./testdata/gc/test1.cache", []byte("value1"), filesMode)
		require.NoError(t, err)

		err = os.WriteFile("./testdata/gc/test2.cache", []byte("value2"), filesMode)
		require.NoError(t, err)

		m1 := newMeta("test1", &ItemOptions{TTL: time.Millisecond}, time.Hour)
		f1, err := os.Create("./testdata/gc/test1.cache--meta")
		require.NoError(t, err)

		err = saveMeta(context.Background(), m1, f1)
		require.NoError(t, err)

		m2 := newMeta("test2", &ItemOptions{TTL: time.Hour}, time.Hour)
		f2, err := os.Create("./testdata/gc/test2.cache--meta")
		require.NoError(t, err)

		err = saveMeta(context.Background(), m2, f2)
		require.NoError(t, err)
	}

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
