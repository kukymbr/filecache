package filecache

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

		err := prepareDir(dir)

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
		err := prepareDir(dir)

		assert.Error(t, err, i)
		assert.NoDirExists(t, dir, i)
	}
}
