package filecache

import "io"

// OpenResult is a result of the file cache's Open operation.
type OpenResult struct {
	hit     bool
	reader  io.ReadCloser
	options *ItemOptions
}

// Hit returns true, if requested key found in cache.
func (r *OpenResult) Hit() bool {
	return r.hit
}

// Reader returns the cached data reader.
func (r *OpenResult) Reader() io.ReadCloser {
	return r.reader
}

// Options returns a found cache item options.
func (r *OpenResult) Options() *ItemOptions {
	return r.options
}

// ReadResult is a result of the file cache's Read operation.
type ReadResult struct {
	hit     bool
	data    []byte
	options *ItemOptions
}

// Hit returns true, if requested key found in cache.
func (r *ReadResult) Hit() bool {
	return r.hit
}

// Data returns the cached data.
func (r *ReadResult) Data() []byte {
	return r.data
}

// Options returns a found cache item options.
func (r *ReadResult) Options() *ItemOptions {
	return r.options
}
