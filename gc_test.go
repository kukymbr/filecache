package filecache

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kukymbr/filecache/v2/internal/util"
	"github.com/stretchr/testify/require"
)

func prepareGCTestFiles(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("./testdata/gc")
	})

	_ = os.RemoveAll("./testdata/gc")

	err := os.MkdirAll("./testdata/gc", util.DirsMode)
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

	// To invalidate test1 item.
	time.Sleep(5 * time.Millisecond)
}
