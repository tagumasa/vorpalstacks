package cloudwatch

import (
	"math"
	"testing"
	"time"
)

// Weighted statistics must equal what a flat one-element-per-sample
// expansion produces. Flat expansion of {[1,3],[2,5],[10,2]} is
// [1,1,1,2,2,2,2,2,10,10].
func TestWeightedStatsMatchFlatExpansion(t *testing.T) {
	var acc sampleAccumulator
	acc.add(1, 3)
	acc.add(2, 5)
	acc.add(10, 2)

	stats := computeExtendedStatistics(acc.samples, []string{
		"p0", "p10", "p50", "p90", "p100", "tm25", "tc1", "ts1", "wm25", "IQM",
	})

	expect := map[string]float64{
		"p0":   1,
		"p10":  1,  // rank 0.9 interpolates within the leading 1s
		"p50":  2,  // rank 4.5 interpolates within the 2s
		"p90":  10, // rank 8.1 interpolates within the trailing 10s
		"p100": 10,
		// trim 2 from each end: [1,2,2,2,2,2] -> (1+2*5)/6
		"tm25": 11.0 / 6.0,
		// trim 1 from each end: [1,1,2,2,2,2,2,10] -> 22/8
		"tc1": 2.75,
		"ts1": 22.0,
		// replace two lowest with sorted[2]=1 (unchanged) and two highest
		// with sorted[7]=2: (33 - 20 + 4)/10
		"wm25": 1.7,
		// mean of sorted[2..6] = [1,2,2,2,2]
		"IQM": 1.8,
	}
	for stat, want := range expect {
		got, ok := stats[stat]
		if !ok {
			t.Fatalf("statistic %q missing from result", stat)
		}
		if math.Abs(got-want) > 1e-9 {
			t.Errorf("%s = %v, want %v", stat, got, want)
		}
	}
}

// A uniform statistic set with a huge SampleCount must resolve analytically
// without allocating per sample; the previous physical expansion allowed a
// client to request gigabytes of memory with a single datum.
func TestWeightedStatsUniformHugeSampleCount(t *testing.T) {
	var acc sampleAccumulator
	acc.add(42.5, 1_000_000_000)
	acc.add(1, 1)

	start := time.Now()
	stats := computeExtendedStatistics(acc.samples, []string{"p0", "p50", "p99.9", "p100"})
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("weighted computation took %v; expansion likely occurred", elapsed)
	}

	for stat, want := range map[string]float64{
		"p0": 1, "p50": 42.5, "p99.9": 42.5, "p100": 42.5,
	} {
		if got := stats[stat]; got != want {
			t.Errorf("%s = %v, want %v", stat, got, want)
		}
	}
}

// Non-positive weights count as a single sample, matching the minimum of
// one underlying data point.
func TestSampleAccumulatorNonPositiveWeightCountsOnce(t *testing.T) {
	var acc sampleAccumulator
	acc.add(7, 0)
	acc.add(9, -3)
	if acc.total != 2 {
		t.Fatalf("total %d, want 2", acc.total)
	}
	stats := computeExtendedStatistics(acc.samples, []string{"p50"})
	if got := stats["p50"]; got != 8 {
		t.Fatalf("p50 = %v, want 8 (midpoint of 7 and 9)", got)
	}
}
