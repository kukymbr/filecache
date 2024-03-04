package filecache_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/kukymbr/filecache/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_WhenValid_ExpectNoError(t *testing.T) {
	tests := []func() (filecache.FileCache, error){
		func() (filecache.FileCache, error) {
			return filecache.New("")
		},
		func() (filecache.FileCache, error) {
			return filecache.NewInTemp()
		},
		func() (filecache.FileCache, error) {
			return filecache.New("./testdata/new")
		},
		func() (filecache.FileCache, error) {
			return filecache.New("./testdata/new", filecache.InstanceOptions{})
		},
		func() (filecache.FileCache, error) {
			return filecache.New("", filecache.InstanceOptions{
				PathGenerator: filecache.FilteredKeyPath,
				DefaultTTL:    time.Hour,
				GCDivisor:     1,
			})
		},
	}

	for i, factory := range tests {
		fc, err := factory()

		assert.NotNil(t, fc, i)
		assert.NoError(t, err, i)
		assert.NotEmpty(t, fc.GetPath())
		assert.NotEqual(t, ".", fc.GetPath())
		assert.DirExists(t, fc.GetPath(), i)
	}
}

func TestNew_WhenInvalid_ExpectError(t *testing.T) {
	tests := []func() (filecache.FileCache, error){
		func() (filecache.FileCache, error) {
			return filecache.New("./testdata/new/item.cache")
		},
		func() (filecache.FileCache, error) {
			return filecache.New("", filecache.InstanceOptions{}, filecache.InstanceOptions{})
		},
	}

	for i, factory := range tests {
		fc, err := factory()

		assert.Nil(t, fc, i)
		assert.Error(t, err, i)
	}
}

func TestFileCache_WriteRead(t *testing.T) {
	target := getTarget(t, "writeread")

	fc, err := filecache.New(target)
	require.NoError(t, err)

	{
		n, err := fc.Write(context.Background(), "test1", strings.NewReader("value1"))

		assert.Equal(t, int64(6), n)
		assert.NoError(t, err)
	}

	{
		n, err := fc.WriteData(
			context.Background(),
			"test2",
			[]byte("value2"),
			filecache.ItemOptions{
				Name:   "Name2",
				TTL:    time.Hour,
				Fields: filecache.NewValues("key1", "val1", "key2", "val2"),
			},
		)

		assert.Equal(t, int64(6), n)
		assert.NoError(t, err)
	}

	{
		res, err := fc.Open(context.Background(), "test1")

		reader := res.Reader()
		options := res.Options()

		data, readErr := io.ReadAll(reader)

		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.True(t, res.Hit())
		assert.NotNil(t, reader)
		assert.NotNil(t, options)
		assert.Equal(t, "value1", string(data))
		assert.NoError(t, readErr)
	}

	{
		res, err := fc.Read(context.Background(), "test2")

		data := res.Data()
		options := res.Options()

		assert.NotNil(t, res)
		assert.NoError(t, err)
		assert.True(t, res.Hit())
		assert.NotNil(t, data)
		assert.NotNil(t, options)
		assert.Equal(t, "value2", string(data))
		assert.Equal(t, "Name2", options.Name)
		assert.Equal(t, time.Hour, options.TTL)
		assert.Equal(t, "val1", options.Fields["key1"])
		assert.Equal(t, "val2", options.Fields["key2"])
	}

	{
		err := fc.Invalidate(context.Background(), "test1")

		assert.NoError(t, err)
	}

	{
		err := fc.Invalidate(context.Background(), "test2")

		assert.NoError(t, err)
	}
}

func TestFileCache_WhenContextCanceled_ExpectError(t *testing.T) {
	fc, err := filecache.New("")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	{
		n, err := fc.Write(ctx, "test1", strings.NewReader("value1"))

		assert.Equal(t, int64(0), n)
		assert.ErrorIs(t, err, context.Canceled)
	}

	{
		n, err := fc.WriteData(ctx, "test2", []byte("value2"))

		assert.Equal(t, int64(0), n)
		assert.ErrorIs(t, err, context.Canceled)
	}

	{
		res, err := fc.Open(ctx, "test1")

		assert.Nil(t, res)
		assert.ErrorIs(t, err, context.Canceled)
	}

	{
		res, err := fc.Read(ctx, "test2")

		assert.Nil(t, res)
		assert.ErrorIs(t, err, context.Canceled)
	}

	{
		err := fc.Invalidate(ctx, "test1")

		assert.ErrorIs(t, err, context.Canceled)
	}
}
