//go:build windows

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
		{`C:/data/cache`, `C:\data\cache`},
		{`C:\data/cache`, `C:\data\cache`},
	}

	for i, test := range tests {
		dir := fixSeparators(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}
