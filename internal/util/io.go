package util

import (
	"context"
	"io"
)

type ReaderFunc func(p []byte) (n int, err error)

func (rf ReaderFunc) Read(p []byte) (n int, err error) { return rf(p) }

// CopyWithCtx is a wrapper for the io.CopyWithCtx function with context handling.
func CopyWithCtx(ctx context.Context, dst io.Writer, src io.Reader) (readN int64, err error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	return io.Copy(dst, ReaderFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
}

// ReadAll is an alias for an io.ReadAll, but with context cancel respect.
func ReadAll(ctx context.Context, r io.Reader) ([]byte, error) {
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
