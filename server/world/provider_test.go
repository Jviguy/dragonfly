package world

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// overlapProvider records whether two of its calls ever ran at once.
type overlapProvider struct {
	NopProvider
	inFlight atomic.Int32
	overlap  atomic.Bool
}

func (p *overlapProvider) SaveSettings(*Settings) {
	if p.inFlight.Add(1) > 1 {
		p.overlap.Store(true)
	}
	time.Sleep(20 * time.Millisecond)
	p.inFlight.Add(-1)
}

// The dimensions of a server share one Provider, which need not be safe for concurrent use, so
// worlds saving at the same moment must still call it one at a time.
func TestWorldsSharingAProviderCallItOneAtATime(t *testing.T) {
	p := &overlapProvider{NopProvider: NopProvider{Set: defaultSettings()}}
	worlds := []*World{
		Config{Provider: p, Dim: Overworld}.New(),
		Config{Provider: p, Dim: Nether}.New(),
		Config{Provider: p, Dim: End}.New(),
	}
	var wg sync.WaitGroup
	for _, w := range worlds {
		wg.Go(w.Save)
	}
	wg.Wait()
	for _, w := range worlds {
		_ = w.Close()
	}
	if p.overlap.Load() {
		t.Fatal("worlds sharing a provider called it concurrently")
	}
}
