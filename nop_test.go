package filecache_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/kukymbrgo/filecache"
)

func TestNopFileCache(t *testing.T) {
	fc := filecache.NewNop()

	{
		n, err := fc.Write(context.Background(), "test", strings.NewReader(""), nil)

		assert.Equal(t, int64(0), n)
		assert.NoError(t, err)
	}

	{
		n, err := fc.WriteData(context.Background(), "test", []byte("test"), nil)

		assert.Equal(t, int64(0), n)
		assert.NoError(t, err)
	}

	{
		r, err := fc.Open(context.Background(), "test")

		assert.NotNil(t, r)
		assert.NoError(t, err)
	}

	{
		data, err := fc.Read(context.Background(), "test")

		assert.NotNil(t, data)
		assert.NoError(t, err)
	}

	{
		err := fc.Invalidate(context.Background(), "test")

		assert.NoError(t, err)
	}
}
