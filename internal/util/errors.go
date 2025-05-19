package util

import "errors"

var (
	ErrDirNotExists = errors.New("directory does not exist")
	ErrNotADir      = errors.New("not a directory")
)
