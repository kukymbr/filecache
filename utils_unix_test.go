//go:build unix

package filecache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixSeparators(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", ""},
		{`/var\cache\fc`, `/var/cache/fc`},
	}

	for i, test := range tests {
		dir := fixSeparators(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}
