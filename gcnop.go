package filecache

type gcNop struct{}

func (g *gcNop) OnInstanceInit() {}

func (g *gcNop) OnOperation() {}

func (g *gcNop) Close() error {
	return nil
}
