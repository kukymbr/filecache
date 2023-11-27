package filecache

import (
	"context"
	"time"
)

type gcInterval struct {
	dir      string
	interval time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	ticker *time.Ticker
}

func (g *gcInterval) OnInstanceInit() {
	g.ticker = time.NewTicker(g.interval)
	g.ctx, g.cancel = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-g.ticker.C:
				g.run()
			case <-g.ctx.Done():
				return
			}
		}
	}()
}

func (g *gcInterval) OnOperation() {}

func (g *gcInterval) Close() error {
	g.ticker.Stop()
	g.cancel()

	return nil
}

func (g *gcInterval) run() {
	scanner := newExpiredScanner(g.dir)

	_ = scanner.Scan(func(entry ScanEntry) error {
		invalidate(entry.itemPath, entry.metaPath)

		return nil
	})
}
