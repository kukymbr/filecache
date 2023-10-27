package filecache

import (
	"io"
	"strings"
)

// NewNop creates no-operation file cache instance.
func NewNop() FileCache {
	return &nopFileCache{}
}

type nopFileCache struct{}

func (fc *nopFileCache) Write(_ string, _ io.Reader, _ *ItemOptions) (written int, err error) {
	return 0, nil
}

func (fc *nopFileCache) WriteData(_ string, _ []byte, _ *ItemOptions) (written int, err error) {
	return 0, nil
}

func (fc *nopFileCache) Open(_ string) (result *OpenResult, err error) {
	return &OpenResult{
		hit: true,
		reader: io.NopCloser(
			strings.NewReader(""),
		),
		options: &ItemOptions{},
	}, nil
}

func (fc *nopFileCache) Read(_ string) (result *ReadResult, err error) {
	return &ReadResult{
		hit:     true,
		data:    []byte(""),
		options: &ItemOptions{},
	}, nil
}
