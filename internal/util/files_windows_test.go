//go:build windows

package util

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
		dir := FixSeparators(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}
