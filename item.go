package filecache

import (
	"os"
)

// Item is a cached item structure
type Item struct {
	File *os.File
	Meta *Meta
	Path string
}

// Close file descriptor
func (i *Item) Close() error {
	return i.File.Close()
}
