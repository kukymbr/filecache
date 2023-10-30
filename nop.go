package filecache

import (
	"context"
	"io"
	"os"
	"strings"
)

// NewNop creates no-operation file cache instance.
func NewNop() FileCache {
	return &nopFileCache{}
}

type nopFileCache struct{}

func (fc *nopFileCache) GetPath() string {
	return os.TempDir()
}

func (fc *nopFileCache) Write(_ context.Context, _ string, _ io.Reader, _ ...ItemOptions) (written int64, err error) {
	return 0, nil
}

func (fc *nopFileCache) WriteData(_ context.Context, _ string, _ []byte, _ ...ItemOptions) (written int64, err error) {
	return 0, nil
}

func (fc *nopFileCache) Open(_ context.Context, _ string) (result *OpenResult, err error) {
	return &OpenResult{
		hit: true,
		reader: io.NopCloser(
			strings.NewReader(""),
		),
		options: &ItemOptions{},
	}, nil
}

func (fc *nopFileCache) Read(_ context.Context, _ string) (result *ReadResult, err error) {
	return &ReadResult{
		hit:     true,
		data:    []byte(""),
		options: &ItemOptions{},
	}, nil
}

func (fc *nopFileCache) Invalidate(_ context.Context, _ string) error {
	return nil
}
