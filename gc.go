package filecache

import (
	"math/rand"
)

func newGarbageCollector(dir string, divisor uint) *garbageCollector {
	return &garbageCollector{
		dir:     dir,
		divisor: divisor,
	}
}

type garbageCollector struct {
	dir     string
	divisor uint
}

func (g *garbageCollector) run() error {
	if !g.decideToRun() {
		return nil
	}

	scanner := newExpiredScanner(g.dir)

	return scanner.Scan(func(entry ScanEntry) error {
		invalidate(entry.itemPath, entry.metaPath)

		return nil
	})
}

func (g *garbageCollector) decideToRun() bool {
	switch g.divisor {
	case 0:
		return false
	case 1:
		return true
	default:
		//nolint:gosec
		i := (rand.Int63n(int64(g.divisor)) + 1) / int64(g.divisor)

		return i == 1
	}
}
