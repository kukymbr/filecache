package filecache

import (
	"math/rand"
	"os"
	"path/filepath"
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
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if gc.fc.fileIsMeta(path) {
			return nil
		}
		if meta := gc.fc.readMeta(path); meta != nil {
			if gc.fc.isExpired(meta) {
				_ = gc.fc.invalidatePath(path)
			}
		}
		return nil
	}
	_ = filepath.Walk(gc.fc.Path(), walkFn)
}
