package filecache_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/kukymbrgo/filecache"
)

func TestPathGenerators(t *testing.T) {
	tests := []struct {
		Key      string
		Fn       filecache.PathGeneratorFn
		Expected string
	}{
		{
			Key:      "test1/test",
			Fn:       filecache.FilteredKeyPath,
			Expected: "test1test",
		},
		{
			Key:      "test2",
			Fn:       filecache.WithExt(filecache.FilteredKeyPath, " .cache "),
			Expected: "test2.cache",
		},
		{
			Key:      "test3",
			Fn:       filecache.WithExt(filecache.HashedKeyPath, ".json"),
			Expected: "3ebfa301dc59196f18593c45e519287a23297589.json",
		},
		{
			Key:      "test4",
			Fn:       filecache.WithExt(filecache.HashedKeySplitPath, "html"),
			Expected: "1f/f2/b3/704aede04eecb51e50ca698efd50a1379b.html",
		},
		{
			Key:      "///",
			Fn:       filecache.FilteredKeyPath,
			Expected: "c64b8480a468e7f8b15abd0cb49a4f9a451af542",
		},
	}

	for i, test := range tests {
		path := test.Fn(test.Key)

		assert.Equal(t, test.Expected, path, i)
	}
}
