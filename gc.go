package filecache

import (
	"math/rand"
	"time"
)

type garbageCollector struct {
	fc *FileCache
}

// execute garbage collector with 1/GCDivisor probability
func (gc *garbageCollector) execute() {
	if !gc.decideToRun() {
		return
	}
	gc.run()
}

func (gc *garbageCollector) decideToRun() bool {
	div := int64(GCDivisor)
	switch div {
	case 0:
		return false
	case 1:
		return true
	default:
		rand.Seed(time.Now().Unix())
		i := (rand.Int63n(div) + 1) / div
		return i == 1
	}
}

func (gc *garbageCollector) run() {
	hitFn := func(meta *Meta, itemPath string, metaPath string) error {
		if meta.IsExpired() {
			_ = invalidatePath(itemPath)
		}
		return nil
	}

	scanner := NewScanner(gc.fc)
	_ = scanner.Scan(hitFn, false, true)
}
