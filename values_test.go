package filecache_test

import (
	"testing"
	"time"

	"github.com/franchb/filecache/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewValues(t *testing.T) {
	var values filecache.Values

	assert.NotPanics(t, func() {
		values = filecache.NewValues(
			"key1", "value1",
			"key2", 2,
			3, 3,
			"key4", time.Second,
		)
	})

	assert.NotNil(t, values)
	assert.Len(t, values, 3)
	assert.Equal(t, "value1", values["key1"])
	assert.Equal(t, 2, values["key2"])
	assert.Equal(t, time.Second, values["key4"])
	assert.NotContains(t, values, "key3")
}
