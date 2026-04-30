package apremind

import (
	"math"
	"math/rand"
	"time"
)

// ApCheckItem represents a user's pending AP check in the priority queue.
type ApCheckItem struct {
	UserNumber    int64
	NextCheckTime time.Time
	heapIndex     int
}

// ApCheckHeap is a min-heap of ApCheckItem, ordered by NextCheckTime (earliest first).
type ApCheckHeap []*ApCheckItem

func (h ApCheckHeap) Len() int { return len(h) }

func (h ApCheckHeap) Less(i, j int) bool {
	return h[i].NextCheckTime.Before(h[j].NextCheckTime)
}

func (h ApCheckHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].heapIndex = i
	h[j].heapIndex = j
}

func (h *ApCheckHeap) Push(x any) {
	item := x.(*ApCheckItem)
	item.heapIndex = len(*h)
	*h = append(*h, item)
}

func (h *ApCheckHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.heapIndex = -1
	*h = old[:n-1]
	return item
}

// betaDelay returns a random delay using a Beta(2, 35) distribution
// scaled to [1, 300] seconds, with the mode around 10 seconds.
func betaDelay() time.Duration {
	x := betaSample(2.0, 35.0)
	seconds := 1.0 + x*299.0
	return time.Duration(seconds * float64(time.Second))
}

// betaSample generates a random variate from Beta(alpha, beta) distribution
// using the Gamma distribution relationship: Beta(a,b) = Gamma(a) / (Gamma(a) + Gamma(b)).
func betaSample(alpha, beta float64) float64 {
	x := gammaSample(alpha)
	y := gammaSample(beta)
	if x+y == 0 {
		return 0.5
	}
	return x / (x + y)
}

// gammaSample generates a random variate from Gamma(shape, 1) distribution
// using the Marsaglia-Tsang method.
func gammaSample(shape float64) float64 {
	if shape < 1 {
		return gammaSample(shape+1) * math.Pow(rand.Float64(), 1.0/shape)
	}
	d := shape - 1.0/3.0
	c := 1.0 / math.Sqrt(9.0*d)
	for {
		var x, v float64
		for {
			x = rand.NormFloat64()
			v = 1.0 + c*x
			if v > 0 {
				break
			}
		}
		v = v * v * v
		u := rand.Float64()
		if u < 1.0-0.0331*(x*x)*(x*x) {
			return d * v
		}
		if math.Log(u) < 0.5*x*x+d*(1.0-v+math.Log(v)) {
			return d * v
		}
	}
}
