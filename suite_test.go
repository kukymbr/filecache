package filecache_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func getTarget(t *testing.T, suiteName string) string {
	target := "./testdata/" + suiteName + "/" + time.Now().Format("20060102150405.99999")

	err := os.MkdirAll(target, 0755)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.RemoveAll(target)
	})

	return target
}
