package apremind

import (
	"container/heap"
	"testing"
	"time"
)

// ── ApCheckHeap ──────────────────────────────────────────────────────────────

// TestApCheckHeapOrdering verifies that heap.Pop returns items in ascending
// NextCheckTime order regardless of insertion order.
func TestApCheckHeapOrdering(t *testing.T) {
	h := make(ApCheckHeap, 0)
	heap.Init(&h)

	now := time.Now()
	// Insert in intentionally non-sorted order.
	for _, pair := range []struct {
		user  int64
		delay time.Duration
	}{
		{3, 3 * time.Minute},
		{1, 1 * time.Minute},
		{2, 2 * time.Minute},
	} {
		heap.Push(&h, &ApCheckItem{
			UserNumber:    pair.user,
			NextCheckTime: now.Add(pair.delay),
		})
	}

	if h.Len() != 3 {
		t.Fatalf("want 3 items, got %d", h.Len())
	}

	wantOrder := []int64{1, 2, 3}
	for _, want := range wantOrder {
		got := heap.Pop(&h).(*ApCheckItem)
		if got.UserNumber != want {
			t.Errorf("want user %d, got %d", want, got.UserNumber)
		}
	}
}

// TestApCheckHeapFix verifies that heap.Fix reorders the heap after
// a key change.
func TestApCheckHeapFix(t *testing.T) {
	h := make(ApCheckHeap, 0)
	heap.Init(&h)

	now := time.Now()
	item1 := &ApCheckItem{UserNumber: 1, NextCheckTime: now.Add(10 * time.Minute)}
	item2 := &ApCheckItem{UserNumber: 2, NextCheckTime: now.Add(20 * time.Minute)}
	heap.Push(&h, item1)
	heap.Push(&h, item2)

	// Move item2 ahead of item1.
	item2.NextCheckTime = now.Add(1 * time.Minute)
	heap.Fix(&h, item2.heapIndex)

	first := heap.Pop(&h).(*ApCheckItem)
	if first.UserNumber != 2 {
		t.Errorf("want user 2 at head after Fix, got %d", first.UserNumber)
	}
}

// TestApCheckHeapRemove verifies that removing an arbitrary item keeps
// the remaining items correctly ordered.
func TestApCheckHeapRemove(t *testing.T) {
	h := make(ApCheckHeap, 0)
	heap.Init(&h)

	now := time.Now()
	items := []*ApCheckItem{
		{UserNumber: 1, NextCheckTime: now.Add(1 * time.Minute)},
		{UserNumber: 2, NextCheckTime: now.Add(2 * time.Minute)},
		{UserNumber: 3, NextCheckTime: now.Add(3 * time.Minute)},
	}
	for _, it := range items {
		heap.Push(&h, it)
	}

	// Remove the middle item.
	heap.Remove(&h, items[1].heapIndex)

	if h.Len() != 2 {
		t.Fatalf("want 2 items after remove, got %d", h.Len())
	}

	first := heap.Pop(&h).(*ApCheckItem)
	second := heap.Pop(&h).(*ApCheckItem)

	if first.UserNumber != 1 || second.UserNumber != 3 {
		t.Errorf("want order [1,3] after remove, got [%d,%d]", first.UserNumber, second.UserNumber)
	}
}

// TestApCheckHeapSwapUpdatesIndex verifies that Swap keeps heapIndex fields
// consistent (critical for heap.Fix and heap.Remove correctness).
func TestApCheckHeapSwapUpdatesIndex(t *testing.T) {
	h := make(ApCheckHeap, 0)
	heap.Init(&h)

	now := time.Now()
	for i := 0; i < 5; i++ {
		heap.Push(&h, &ApCheckItem{
			UserNumber:    int64(i + 1),
			NextCheckTime: now.Add(time.Duration(5-i) * time.Minute),
		})
	}

	// After building the heap every item's heapIndex must match its actual
	// position in the underlying slice.
	for i, item := range h {
		if item.heapIndex != i {
			t.Errorf("item %d: heapIndex=%d, want %d", item.UserNumber, item.heapIndex, i)
		}
	}
}

// ── Beta / Gamma distribution ────────────────────────────────────────────────

// TestBetaDelayRange verifies that betaDelay always produces a value in [1s, 300s].
func TestBetaDelayRange(t *testing.T) {
	const iterations = 2000
	for i := 0; i < iterations; i++ {
		d := betaDelay()
		if d < time.Second || d > 300*time.Second {
			t.Errorf("betaDelay() = %v, want in [1s, 300s]", d)
		}
	}
}

// TestBetaSampleRange verifies that betaSample always produces a value in [0, 1].
func TestBetaSampleRange(t *testing.T) {
	const iterations = 2000
	for i := 0; i < iterations; i++ {
		x := betaSample(2.0, 35.0)
		if x < 0 || x > 1 {
			t.Errorf("betaSample() = %v, want in [0, 1]", x)
		}
	}
}

// TestGammaSamplePositive verifies that gammaSample always returns a positive
// value across a range of shape parameters including the shape < 1 branch.
func TestGammaSamplePositive(t *testing.T) {
	shapes := []float64{0.3, 0.5, 1.0, 2.0, 5.0, 35.0}
	const iterations = 200
	for _, shape := range shapes {
		for i := 0; i < iterations; i++ {
			x := gammaSample(shape)
			if x <= 0 {
				t.Errorf("gammaSample(%v) = %v, want > 0", shape, x)
			}
		}
	}
}

// TestBetaDelayDistribution verifies that the Beta(2,35) distribution
// concentrates most of its mass near the lower end of the [1,300] range.
// Because the mode of Beta(2,35) is ~(2-1)/(2+35-2) = 1/35 ≈ 3% → ~10s,
// we expect well over 60% of samples to fall below 30 s.
func TestBetaDelayDistribution(t *testing.T) {
	const iterations = 3000
	below30 := 0
	for i := 0; i < iterations; i++ {
		if betaDelay() < 30*time.Second {
			below30++
		}
	}
	// Conservative threshold: at least 60 % below 30 s.
	threshold := iterations * 60 / 100
	if below30 < threshold {
		t.Errorf("want ≥%d samples below 30s (≥60%%), got %d (%.1f%%)",
			threshold, below30, float64(below30)*100/float64(iterations))
	}
}
