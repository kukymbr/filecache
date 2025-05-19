package filecache

import (
	"math/rand"

	"github.com/kukymbr/filecache/v2/internal/util"
)

type gcProbability struct {
	dir string

	onInitDivisor uint
	onOpDivisor   uint
}

func (g *gcProbability) OnInstanceInit() {
	g.run(g.onInitDivisor)
}

func (g *gcProbability) OnOperation() {
	g.run(g.onOpDivisor)
}

func (g *gcProbability) Close() error {
	return nil
}

func (g *gcProbability) run(divisor uint) {
	if !g.decideToRun(divisor) {
		return
	}

	scanner := newExpiredScanner(g.dir)

	_ = scanner.Scan(func(entry ScanEntry) error {
		util.DeleteCacheFiles(entry.itemPath, entry.metaPath)

		return nil
	})
}

func (g *gcProbability) decideToRun(divisor uint) bool {
	switch divisor {
	case 0:
		return false
	case 1:
		return true
	default:
		//nolint:gosec
		i := (rand.Int63n(int64(divisor)) + 1) / int64(divisor)

		return i == 1
	}
}
