package filecache_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/kukymbrgo/filecache"
)

func TestScanner(t *testing.T) {
	fc, err := filecache.New(getTarget(t, "scanner"))
	require.NoError(t, err)

	// Prepare cache
	{
		_, err = fc.WriteData(context.Background(), "test1", []byte("value1"))
		require.NoError(t, err)

		_, err = fc.WriteData(
			context.Background(),
			"test2",
			[]byte("value2"),
			filecache.ItemOptions{TTL: time.Hour},
		)
		require.NoError(t, err)

		_, err = fc.WriteData(
			context.Background(),
			"test3",
			[]byte("value3"),
			filecache.ItemOptions{TTL: time.Millisecond},
		)
		require.NoError(t, err)
	}

	time.Sleep(2 * time.Millisecond)

	scanner := filecache.NewScanner(fc.GetPath())
	scannedKeys := make([]string, 0)

	err = scanner.Scan(func(entry filecache.ScanEntry) error {
		scannedKeys = append(scannedKeys, entry.Key)

		return nil
	})

	assert.Len(t, scannedKeys, 2)
	assert.Contains(t, scannedKeys, "test1")
	assert.Contains(t, scannedKeys, "test2")
	assert.NotContains(t, scannedKeys, "test3")
}
