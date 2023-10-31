package filecache

import (
	"context"
	"fmt"
	"io"
	"os"
)

const filesMode os.FileMode = 0644

func create(key string, path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filesMode)
	if err != nil {
		return nil, fmt.Errorf("failed to create file for cache key %s: %w", key, err)
	}

	return f, nil
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

// copyWithCtx is a wrapper for the io.copyWithCtx function with context handling.
func copyWithCtx(ctx context.Context, dst io.Writer, src io.Reader) (readN int64, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return io.Copy(dst, readerFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
}

// readAll is an alias for an io.ReadAll, but with context cancel respect.
func readAll(ctx context.Context, r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}

		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]

		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return b, err
		}
	}
}
