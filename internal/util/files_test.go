package util

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateDir_WhenValid_ExpectNoError(t *testing.T) {
	dirs := []string{
		".",
		"./testdata",
		"./testdata/utils",
		os.TempDir(),
	}

	for i, dir := range dirs {
		err := validateDir(dir)

		assert.NoError(t, err, i)
		assert.DirExists(t, dir, i)
	}
}

func TestValidateDir_WhenInvalid_ExpectError(t *testing.T) {
	dirs := []string{
		"./testdata/utils/unknown",
		"./testdata/utils/test.txt",
	}

	for i, dir := range dirs {
		err := validateDir(dir)

		assert.Error(t, err, i)
		assert.NoDirExists(t, dir, i)
	}
}

func TestPrepareDir_WhenValid_ExpectNoError(t *testing.T) {
	dirs := []string{
		"./testdata/utils",
		"./testdata/utils/test_" + fmt.Sprintf("%v", time.Now().UnixNano()),
		os.TempDir() + "/test_" + fmt.Sprintf("%v", time.Now().UnixNano()),
	}

	for i, dir := range dirs {
		dir := dir
		existed := validateDir(dir) == nil

		err := PrepareDir(dir)

		assert.NoError(t, err, i)
		assert.DirExists(t, dir, i)

		if !existed {
			t.Cleanup(func() {
				_ = os.RemoveAll(dir)
			})
		}
	}
}

func TestPrepareDir_WhenInvalid_ExpectError(t *testing.T) {
	dirs := []string{
		"./testdata/utils/test.txt",
	}

	for i, dir := range dirs {
		err := PrepareDir(dir)

		assert.Error(t, err, i)
		assert.NoDirExists(t, dir, i)
	}
}

func TestItemFilesValid_WhenValid_ExpectTrue(t *testing.T) {
	tests := []struct {
		Item string
		Meta string
	}{
		{"./testdata/utils/cache/item.cache", "./testdata/utils/cache/item.cache--meta"},
	}

	for i, test := range tests {
		res := ItemFilesValid(test.Item, test.Meta)

		assert.True(t, res, i)
	}
}

func TestItemFilesValid_WhenInvalid_ExpectFalse(t *testing.T) {
	tests := []struct {
		Item string
		Meta string
	}{
		{"", ""},
		{"./testdata/utils/cache/item.cache", ""},
		{"./testdata/utils/cache/item.cache", "./testdata/utils/cache"},
		{"./testdata/utils/cache/item.cache", "./testdata/utils/cache/unknown--meta"},
		{"", "./testdata/utils/cache/item.cache--meta"},
		{"./testdata/utils/cache", "./testdata/utils/cache/item.cache--meta"},
		{"./testdata/utils/cache/unknown--meta", "./testdata/utils/cache/item.cache--meta"},
	}

	for i, test := range tests {
		res := ItemFilesValid(test.Item, test.Meta)

		assert.False(t, res, i)
	}
}

func TestIsExpired(t *testing.T) {
	tests := []struct {
		CreatedAt time.Time
		TTL       time.Duration
		Expected  bool
	}{
		{
			CreatedAt: time.Now().Add(-time.Hour * 2),
			TTL:       time.Hour * 3,
			Expected:  false,
		},
		{
			CreatedAt: time.Now().Add(-time.Hour * 2),
			TTL:       TTLEternal,
			Expected:  false,
		},
		{
			CreatedAt: time.Now().Add(-time.Hour * 2),
			TTL:       0,
			Expected:  false,
		},
		{
			CreatedAt: time.Now().Add(-time.Hour * 2),
			TTL:       time.Hour,
			Expected:  true,
		},
	}

	for i, test := range tests {
		expired := IsExpired(test.CreatedAt, test.TTL)

		assert.Equal(t, test.Expected, expired, i)
	}
}
