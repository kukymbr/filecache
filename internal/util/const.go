package util

import (
	"os"
	"time"
)

const (
	// TTLEternal is a TTL value for eternal cache.
	TTLEternal = time.Duration(-1)

	MetaSuffix             = "--meta"
	DirsMode   os.FileMode = 0755
)
