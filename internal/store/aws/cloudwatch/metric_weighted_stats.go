package cloudwatch

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

// weightedSample is one distinct metric value together with the number of
// underlying samples it represents.
type weightedSample struct {
	value  float64
	weight int
}

// sampleAccumulator collects period samples as weighted pairs. Statistic
// sets with Min == Max and multi-value entries with Counts represent many
// identical samples with a single number; physically expanding them to one
// slice element per sample is a memory-amplification vector because the
// client-supplied SampleCount or Count is unbounded (a SampleCount of 1e9
// would allocate roughly 8 GB before any percentile is computed).
type sampleAccumulator struct {
	samples []weightedSample
	total   int
}

// add records weight samples of value. A non-positive weight counts as a
// single sample, matching the minimum of one underlying data point.
func (a *sampleAccumulator) add(value float64, weight int) {
	if weight <= 0 {
		weight = 1
	}
	a.samples = append(a.samples, weightedSample{value: value, weight: weight})
	a.total += weight
}

// weightedStats indexes accumulated samples for order-statistic queries.
// Every statistic computed over the weighted samples equals the value a
// flat one-element-per-sample expansion would produce, because a block of
// w identical values is a contiguous run in the sorted order.
type weightedStats struct {
	samples []weightedSample // sorted by value ascending
	prefix  []int            // prefix[i] = flat index where block i starts
	n       int              // total number of samples
	sum     float64          // sum of all samples
}

func newWeightedStats(samples []weightedSample) *weightedStats {
	ws := &weightedStats{samples: samples}
	sort.Slice(ws.samples, func(i, j int) bool { return ws.samples[i].value < ws.samples[j].value })
	idx := 0
	for i := range ws.samples {
		ws.prefix = append(ws.prefix, idx)
		idx += ws.samples[i].weight
		ws.sum += ws.samples[i].value * float64(ws.samples[i].weight)
	}
	ws.n = idx
	return ws
}

// valueAt returns the value of the flat sample index i.
func (ws *weightedStats) valueAt(i int) float64 {
	b := sort.Search(len(ws.samples), func(k int) bool {
		return ws.prefix[k]+ws.samples[k].weight > i
	})
	if b >= len(ws.samples) {
		b = len(ws.samples) - 1
	}
	return ws.samples[b].value
}

// rangeSum returns the sum of flat sample indices [from, to).
func (ws *weightedStats) rangeSum(from, to int) float64 {
	var s float64
	for b := range ws.samples {
		lo := ws.prefix[b]
		hi := lo + ws.samples[b].weight
		if hi <= from || lo >= to {
			continue
		}
		l := lo
		if from > l {
			l = from
		}
		h := hi
		if to < h {
			h = to
		}
		s += ws.samples[b].value * float64(h-l)
	}
	return s
}

// rangeMean returns the mean of flat sample indices [from, to).
func (ws *weightedStats) rangeMean(from, to int) float64 {
	if to <= from {
		return 0
	}
	return ws.rangeSum(from, to) / float64(to-from)
}

// percentile interpolates the p-quantile (p in [0, 1]) exactly as a flat
// sorted expansion would.
func (ws *weightedStats) percentile(p float64) float64 {
	if ws.n == 0 {
		return 0
	}
	if ws.n == 1 {
		return ws.samples[0].value
	}
	rank := p * float64(ws.n-1)
	lo := int(math.Floor(rank))
	hi := int(math.Ceil(rank))
	vLo := ws.valueAt(lo)
	vHi := ws.valueAt(hi)
	return vLo + (rank-float64(lo))*(vHi-vLo)
}

// trimmedMean returns the mean after trimming a fraction of samples from
// each end.
func (ws *weightedStats) trimmedMean(trimFrac float64) float64 {
	trimCount := int(float64(ws.n) * trimFrac)
	if trimCount*2 >= ws.n {
		return 0
	}
	return ws.rangeMean(trimCount, ws.n-trimCount)
}

// winsorizedMean replaces the trimmed tail samples with the boundary values
// and returns the resulting mean.
func (ws *weightedStats) winsorizedMean(trimFrac float64) float64 {
	trimCount := int(float64(ws.n) * trimFrac)
	if trimCount*2 >= ws.n {
		return ws.sum / float64(ws.n)
	}
	low := ws.valueAt(trimCount)
	high := ws.valueAt(ws.n - 1 - trimCount)
	s := ws.sum
	s -= ws.rangeSum(0, trimCount)
	s += float64(trimCount) * low
	s -= ws.rangeSum(ws.n-trimCount, ws.n)
	s += float64(trimCount) * high
	return s / float64(ws.n)
}

// interquartileMean returns the mean of the samples between the first and
// third quartile indices.
func (ws *weightedStats) interquartileMean() float64 {
	q1 := int(0.25 * float64(ws.n-1))
	q3 := int(0.75 * float64(ws.n-1))
	if q3 <= q1 {
		return ws.sum / float64(ws.n)
	}
	return ws.rangeMean(q1, q3+1)
}

// computeExtendedStatistics calculates extended statistics (percentiles,
// IQM, trimmed/winsorized means) from the accumulated weighted samples.
func computeExtendedStatistics(samples []weightedSample, stats []string) map[string]float64 {
	ws := newWeightedStats(samples)
	result := make(map[string]float64, len(stats))
	for _, stat := range stats {
		stat = strings.TrimSpace(stat)
		lower := strings.ToLower(stat)
		val, ok := ws.computeExtendedStatistic(stat, lower)
		if ok {
			result[stat] = val
		}
	}
	return result
}

// computeExtendedStatistic computes a single extended statistic. Returns
// (value, true) on success.
func (ws *weightedStats) computeExtendedStatistic(stat, lower string) (float64, bool) {
	if ws.n == 0 {
		return 0, false
	}
	// Percentile: p{n} (e.g. p90, p99.9)
	if strings.HasPrefix(lower, "p") {
		percentile, err := strconv.ParseFloat(stat[1:], 64)
		if err != nil || percentile < 0 || percentile > 100 {
			return 0, false
		}
		return ws.percentile(percentile / 100.0), true
	}
	// Trimmed mean: tm{n}
	if strings.HasPrefix(lower, "tm") {
		percentile, err := strconv.ParseFloat(stat[2:], 64)
		if err != nil || percentile < 0 || percentile >= 50 {
			return 0, false
		}
		return ws.trimmedMean(percentile / 100.0), true
	}
	// Trimmed count: tc{n}
	if strings.HasPrefix(lower, "tc") {
		n, err := strconv.Atoi(stat[2:])
		if err != nil || n < 0 || n*2 >= ws.n {
			return 0, false
		}
		return ws.rangeMean(n, ws.n-n), true
	}
	// Trimmed sum: ts{n}
	if strings.HasPrefix(lower, "ts") {
		n, err := strconv.Atoi(stat[2:])
		if err != nil || n < 0 || n*2 >= ws.n {
			return 0, false
		}
		return ws.rangeSum(n, ws.n-n), true
	}
	// Winsorized mean: wm{n}
	if strings.HasPrefix(lower, "wm") {
		percentile, err := strconv.ParseFloat(stat[2:], 64)
		if err != nil || percentile < 0 || percentile >= 50 {
			return 0, false
		}
		return ws.winsorizedMean(percentile / 100.0), true
	}
	// Interquartile mean: IQM
	if lower == "iqm" {
		return ws.interquartileMean(), true
	}
	return 0, false
}
